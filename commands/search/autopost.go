package search

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search"
)

func (c *commands) autopost(ctx bcr.Contexter) (err error) {
	ch, err := ctx.GetChannelFlag("channel")
	if v, ok := ctx.(*bcr.Context); ok {
		ch, err = v.ParseChannel(v.Args[0])
	}
	if err != nil || ch.GuildID != ctx.GetChannel().GuildID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews) {
		return ctx.SendEphemeral(":x: Couldn't find that channel, or it's not in this server.")
	}

	dur, err := durationparser.Parse(ctx.GetStringFlag("interval"))
	if v, ok := ctx.(*bcr.Context); ok {
		dur, err = durationparser.Parse(v.Args[1])
	}

	if err != nil || (dur < 2*time.Hour && dur != 0) || dur > 7*24*time.Hour {
		s := ctx.GetStringFlag("interval")
		if v, ok := ctx.(*bcr.Context); ok {
			s = v.Args[1]
		}
		if s == "clear" || s == "0" || s == "disable" || s == "off" || s == "reset" {
			dur = 0
		} else {
			return ctx.SendEphemeral(":x: Couldn't parse ``" + bcr.EscapeBackticks(ctx.GetStringFlag("interval")) + "`` as a valid duration (minimum 2 hours, maximum 1 week).")
		}
	}

	var catID *int
	if s := ctx.GetStringFlag("category"); s != "" {
		category, err := c.DB.CategoryID(s)
		if err != nil {
			return ctx.SendEphemeral("Couldn't find a category with the name ``" + bcr.EscapeBackticks(s) + "``.")
		}
		catID = &category
	}

	perms := discord.CalcOverwrites(*ctx.GetGuild(), *ch, *ctx.GetMember())
	if !perms.Has(discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionManageRoles) {
		return ctx.SendEphemeral("You don't have permissions to set autoposts in " + ch.Mention() + ".")
	}

	var count int
	err = c.DB.QueryRow(context.Background(), "select count(*) from autopost where guild_id = $1", ctx.GetGuild().ID.String()).Scan(&count)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	if count >= 10 {
		var exists bool
		err = c.DB.QueryRow(context.Background(), "select exists(select * from autopost where channel_id = $1)", ch.ID).Scan(&exists)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
		if !exists {
			return ctx.SendfX("Only 10 channels can be autoposted to at the same time. Remove one to add a new channel.")
		}
	}

	err = c.setAutopost(ctx.GetGuild().ID, ch.ID, catID, dur)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	if dur == 0 {
		return ctx.SendfX("Disabled autoposting in %v!", ch.Mention())
	}

	s := fmt.Sprintf("Autopost interval set! A random term")
	if ctx.GetStringFlag("category") != "" {
		s += fmt.Sprintf(" from the %v category", ctx.GetStringFlag("category"))
	}
	s += fmt.Sprintf(" will be posted to %v every %v.", ch.Mention(), bcr.HumanizeDuration(bcr.DurationPrecisionMinutes, dur))

	return ctx.SendX(s)
}

func (c *commands) setAutopost(guildID discord.GuildID, channelID discord.ChannelID, categoryID *int, duration time.Duration) (err error) {
	if duration == 0 {
		db.Debug("Removing autopost in %v", channelID)
		_, err = c.DB.Exec(context.Background(), "delete from autopost where channel_id = $1", channelID)
		return
	}

	db.Debug("Setting autopost in %v", channelID)
	_, err = c.DB.Exec(context.Background(), `insert into autopost (guild_id, channel_id, category_id, next_post, interval)
values ($1, $2, $3, $4, $5)
on conflict (channel_id) do update
set category_id = $3, next_post = $4, interval = $5`, guildID.String(), channelID, categoryID, time.Now().UTC().Add(duration), duration)
	return
}

func (c *commands) setNextTime(channelID discord.ChannelID, time time.Time) (err error) {
	db.Debug("Updating next post time in %v to %s", channelID, time)
	_, err = c.DB.Exec(context.Background(), "update autopost set next_post = $1 where channel_id = $2", time, channelID)
	return
}

// Autopost ...
type Autopost struct {
	GuildID   string
	ChannelID discord.ChannelID
	NextPost  time.Time
	Interval  time.Duration

	CategoryID *int
}

func (c *commands) autopostLoop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-sc:
			break
		default:
			aps := []Autopost{}

			err := pgxscan.Select(context.Background(), c.DB, &aps, "select * from autopost where next_post < $1 limit 5", time.Now().UTC())
			if err != nil {
				c.Sugar.Errorf("Error getting autopost info: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, ap := range aps {
				err = c.doAutopost(ap)
				if err != nil {
					c.Sugar.Errorf("Error running autopost in %v: %v", ap.ChannelID, err)
				}
			}
		}
	}
}

func (c *commands) doAutopost(ap Autopost) (err error) {
	sf, _ := discord.ParseSnowflake(ap.GuildID)
	s, _ := c.Router.StateFromGuildID(discord.GuildID(sf))

	perms, err := s.Permissions(ap.ChannelID, c.Router.Bot.ID)
	if err != nil {
		c.setAutopost(0, ap.ChannelID, nil, 0)
		return errors.Wrap(err, "getting permissions")
	}

	if !perms.Has(discord.PermissionViewChannel | discord.PermissionSendMessages) {
		c.Sugar.Errorf("Can't send messages in %v (guild %v), disabling autopost there.", ap.ChannelID, ap.GuildID)
		return c.setAutopost(0, ap.ChannelID, nil, 0)
	}

	var t *search.Term
	if ap.CategoryID != nil {
		t, err = c.DB.RandomTermCategory(*ap.CategoryID, []string{})
	} else {
		t, err = c.DB.RandomTerm([]string{})
	}
	if err != nil {
		return errors.Wrap(err, "get random term")
	}

	_, err = s.SendEmbeds(ap.ChannelID, c.DB.TermEmbed(t))
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	toAdd := ap.Interval
	secs := rand.Intn(60)
	if rand.Intn(2) == 0 {
		secs -= 60
	}
	toAdd += time.Duration(secs) * time.Second

	err = c.setNextTime(ap.ChannelID, time.Now().UTC().Add(toAdd))
	if err != nil {
		return errors.Wrap(err, "set next post time")
	}
	return nil
}

func (c *commands) autopostList(ctx *bcr.Context) (err error) {
	aps := []Autopost{}
	err = pgxscan.Select(context.Background(), c.DB, &aps, "select * from autopost where guild_id = $1 order by channel_id", ctx.Message.GuildID.String())
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	if len(aps) == 0 {
		_, err = ctx.Reply("No channels in this server are being autoposted to.")
		return
	}

	s := ""

	for _, ap := range aps {
		s += fmt.Sprintf("%v: every %v (next post: <t:%v>)", ap.ChannelID.Mention(), bcr.HumanizeDuration(bcr.DurationPrecisionMinutes, ap.Interval), ap.NextPost.Unix())
		if ap.CategoryID != nil {
			if cat := c.DB.CategoryFromID(*ap.CategoryID); cat != nil {
				s += "\nPosting terms from the " + cat.Name + " category"
			}
		}
		s += "\n\n"
	}

	e := discord.Embed{
		Title:       "Autopost channels for " + ctx.Guild.Name,
		Color:       db.EmbedColour,
		Description: strings.TrimSpace(s),
	}
	_, err = ctx.Send("", e)
	return
}
