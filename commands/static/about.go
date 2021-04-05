package static

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"runtime"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/dustin/go-humanize"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

// uh this isn't semver, this is just. increment the number when we feel like it i guess
const botVersion = "v7"

var gitVer string

func init() {
	git := exec.Command("git", "rev-parse", "--short", "HEAD")
	// ignoring errors *should* be fine? if there's no output we just fall back to "unknown"
	b, _ := git.Output()
	gitVer = string(b)
	if gitVer == "" {
		gitVer = "unknown"
	}
}

type category struct {
	ID    int
	Name  string
	Count int
}

func (c *Commands) about(ctx *bcr.Context) (err error) {
	// get term count
	var (
		total = c.DB.TermCount()

		pronouns   int
		categories []category
	)
	err = pgxscan.Select(context.Background(), c.DB.Pool, &categories, "select categories.id, categories.name, count(terms.id) from categories inner join terms on categories.id = terms.category group by categories.id order by categories.id")
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	err = c.DB.Pool.QueryRow(context.Background(), "select count(id) from pronouns").Scan(&pronouns)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	terms := discord.EmbedField{
		Name:   "Terms",
		Value:  fmt.Sprintf("**%v** total", total),
		Inline: true,
	}
	for _, c := range categories {
		terms.Value += fmt.Sprintf("\n**%v** %v terms", c.Count, c.Name)
	}
	terms.Value += fmt.Sprintf("\n\n**%v** pronouns", pronouns)

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	fields := []discord.EmbedField{
		{
			Name:   "Bot version",
			Value:  fmt.Sprintf("%v-%v ([bcr](https://github.com/starshine-sys/bcr) v%v)", botVersion, gitVer, bcr.Version()),
			Inline: true,
		},
		{
			Name:   "Go version",
			Value:  runtime.Version(),
			Inline: true,
		},
		{
			Name:   "Invite",
			Value:  fmt.Sprintf("[Invite link](%v)", invite(ctx)),
			Inline: true,
		},
	}

	if c.Config.Sharded {
		fields = append(fields, discord.EmbedField{
			Name:   "Shard",
			Value:  fmt.Sprintf("#%v (%v total)", c.Router.State.Gateway.Identifier.Shard.ShardID(), c.Router.State.Gateway.Identifier.Shard.NumShards()),
			Inline: true,
		})
	}

	fields = append(fields, []discord.EmbedField{
		{
			Name: "Uptime",
			Value: fmt.Sprintf(
				"%v\n(Since %v)",
				prettyDurationString(time.Since(c.start)),
				c.start.Format("Jan _2 2006, 15:04:05 MST"),
			),
			Inline: true,
		},
		{
			Name:   "Memory used",
			Value:  fmt.Sprintf("%v / %v (%v garbage collected)\n%v goroutines", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc), runtime.NumGoroutine()),
			Inline: false,
		},
		terms,
		{
			Name:   "Source code",
			Value:  fmt.Sprintf("[GitHub](%v)\n/ [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html) license", c.Config.Bot.Git),
			Inline: true,
		},
	}...)

	embed := &discord.Embed{
		Title: "About",
		Color: db.EmbedColour,
		Footer: &discord.EmbedFooter{
			Text: "Made with Arikawa",
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: ctx.Bot.AvatarURL(),
		},
		Timestamp: discord.NewTimestamp(time.Now()),
		Fields:    fields,
	}

	_, err = ctx.Send("", embed)
	return
}

func invite(ctx *bcr.Context) string {
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

func prettyDurationString(duration time.Duration) (out string) {
	var days, hours, hoursFrac, minutes float64

	hours = duration.Hours()
	hours, hoursFrac = math.Modf(hours)
	minutes = hoursFrac * 60

	hoursFrac = math.Mod(hours, 24)
	days = (hours - hoursFrac) / 24
	hours = hours - (days * 24)
	minutes = minutes - math.Mod(minutes, 1)

	if days != 0 {
		out += fmt.Sprintf("%v days, ", days)
	}
	if hours != 0 {
		out += fmt.Sprintf("%v hours, ", hours)
	}
	out += fmt.Sprintf("%v minutes", minutes)

	return
}
