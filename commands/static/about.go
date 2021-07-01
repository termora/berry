package static

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/dustin/go-humanize"
	"github.com/dustin/go-humanize/english"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

var gitVer string

func init() {
	git := exec.Command("git", "rev-parse", "--short", "HEAD")
	// ignoring errors *should* be fine? if there's no output we just fall back to "unknown"
	b, _ := git.Output()
	gitVer = strings.TrimSpace(string(b))
	if gitVer == "" {
		gitVer = "[unknown]"
	}
}

type category struct {
	ID    int
	Name  string
	Count int
}

func (c *Commands) about(ctx *bcr.Context) (err error) {
	t := time.Now()

	msg, err := ctx.Send("...", nil)
	if err != nil {
		return err
	}

	latency := time.Since(t).Round(time.Millisecond)

	// this will return 0ms in the first minute after the bot is restarted
	// can't do much about that though
	heartbeat := ctx.State.Gateway.PacerLoop.EchoBeat.Time().Sub(ctx.State.Gateway.PacerLoop.SentBeat.Time()).Round(time.Millisecond)

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Bot.AvatarURL(),
			Name: "About " + ctx.Bot.Username,
		},
		Color: db.EmbedColour,
		Fields: []discord.EmbedField{
			{
				Name:   "Ping",
				Value:  fmt.Sprintf("Heartbeat: %v\nMessage: %v", heartbeat, latency),
				Inline: true,
			},
			{
				Name:   "Go version",
				Value:  fmt.Sprintf("%v\n%v/%v", runtime.Version(), runtime.GOOS, runtime.GOARCH),
				Inline: true,
			},
		},
		Footer: &discord.EmbedFooter{
			Text: "Version " + gitVer,
		},
		Timestamp: discord.NowTimestamp(),
	}

	if c.Config.Sharded {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:   "Shard",
			Value:  fmt.Sprintf("#%v (%v total)", c.Router.State.Gateway.Identifier.Shard.ShardID(), c.Router.State.Gateway.Identifier.Shard.NumShards()),
			Inline: true,
		})
	}

	servers := humanize.Comma(c.GuildCount)
	if c.Config.Sharded {
		servers += fmt.Sprintf("\nShard %v of %v", c.Router.State.Gateway.Identifier.Shard.ShardID()+1, c.Router.State.Gateway.Identifier.Shard.NumShards())
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name:   "Servers",
		Value:  servers,
		Inline: true,
	}, discord.EmbedField{
		Name: "Memory usage",
		Value: fmt.Sprintf(
			"%v / %v (%v garbage collected)\n%v goroutines",
			humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys),
			humanize.Bytes(stats.TotalAlloc), runtime.NumGoroutine(),
		),
		Inline: true,
	}, discord.EmbedField{
		Name:   "Uptime",
		Value:  fmt.Sprintf("%v\nSince %v UTC", bcr.HumanizeDuration(bcr.DurationPrecisionSeconds, time.Since(c.start)), c.start.Format("2006-01-02 15:04:05")),
		Inline: true,
	})

	// get term count
	var (
		total = c.DB.TermCount()

		pronouns   int
		categories []category
	)

	con, cancel := c.DB.Context()
	defer cancel()

	err = pgxscan.Select(con, c.DB.Pool, &categories, `select
	categories.id, categories.name, count(terms.id)
	from categories
	inner join terms on categories.id = terms.category
	group by categories.id order by categories.id`)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	con, cancel = c.DB.Context()
	defer cancel()

	err = c.DB.Pool.QueryRow(con, "select count(id) from pronouns").Scan(&pronouns)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	{
		slice := []string{}

		for _, c := range categories {
			slice = append(slice, fmt.Sprintf("%v %v terms", c.Count, c.Name))
		}

		e.Fields = append(e.Fields, discord.EmbedField{
			Name: "Numbers",
			Value: fmt.Sprintf(
				"%v terms (%v)\n%v pronouns",
				total,
				english.OxfordWordSeries(slice, "and"),
				pronouns,
			),
		}, discord.EmbedField{
			Name:   "Invite",
			Value:  c.invite(ctx),
			Inline: true,
		}, discord.EmbedField{
			Name:   "Source code",
			Value:  fmt.Sprintf("[GitHub](%v) / [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html) license", c.Config.Bot.Git),
			Inline: true,
		})
	}

	_, err = ctx.Edit(msg, "", &e)
	return
}

func (c *Commands) invite(ctx *bcr.Context) string {
	if c.Config.Bot.CustomInvite != "" {
		return c.Config.Bot.CustomInvite
	}

	// perms is the list of permissions the bot will be granted by default
	var perms = discord.PermissionViewChannel +
		discord.PermissionReadMessageHistory +
		discord.PermissionSendMessages +
		discord.PermissionManageMessages +
		discord.PermissionEmbedLinks +
		discord.PermissionAttachFiles +
		discord.PermissionUseExternalEmojis +
		discord.PermissionAddReactions

	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&scope=applications.commands%%20bot", ctx.Bot.ID, perms)
}

func urlParse(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}
