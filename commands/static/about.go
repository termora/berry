package static

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
)

var botVersion = "v0.3"

func (c *Commands) about(ctx *bcr.Context) (err error) {
	c.cmdMutex.RLock()
	defer c.cmdMutex.RUnlock()
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
		Fields: []discord.EmbedField{
			{
				Name:   "Bot version",
				Value:  fmt.Sprintf("%v (bcr v%v)", botVersion, bcr.Version()),
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
			{
				Name: "Uptime",
				Value: fmt.Sprintf(`%v
				(Since %v)
				
				**Commands run since last restart:** %v (%.1f/min)

				**Terms:** %v
				**Searches since last restart:** %v`,
					prettyDurationString(time.Since(c.start)),
					c.start.Format("Jan _2 2006, 15:04:05 MST"),
					c.cmdCount,
					float64(c.cmdCount)/time.Since(c.start).Minutes(),
					c.db.TermCount(),
					db.GetCount(),
				),
				Inline: false,
			},
			{
				Name:   "Credits",
				Value:  fmt.Sprintf("Check `%vcredits`!", ctx.Router.Prefixes[0]),
				Inline: true,
			},
			{
				Name:   "Source code",
				Value:  "[GitHub](https://github.com/Starshine113/berry)\n/ Licensed under the [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html)",
				Inline: true,
			},
		},
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

	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&scope=bot", ctx.Bot.ID, perms)
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
