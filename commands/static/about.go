package static

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
)

var botVersion = "v0.2"

func (c *Commands) about(ctx *crouter.Ctx) (err error) {
	owner, err := ctx.Session.User(c.config.Bot.BotOwners[0])
	if err != nil {
		return ctx.CommandError(err)
	}

	c.cmdMutex.RLock()
	defer c.cmdMutex.RUnlock()
	embed := &discordgo.MessageEmbed{
		Title: "About",
		Color: db.EmbedColour,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Made with discordgo %v", discordgo.VERSION),
			IconURL: "https://raw.githubusercontent.com/bwmarrin/discordgo/master/docs/img/discordgo.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ctx.BotUser.AvatarURL("256"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Bot version",
				Value:  fmt.Sprintf("%v (crouter v%v)", botVersion, crouter.Version()),
				Inline: true,
			},
			{
				Name:   "discordgo version",
				Value:  fmt.Sprintf("%v (%v)", discordgo.VERSION, runtime.Version()),
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
				Name:   "Author",
				Value:  owner.Mention() + " / " + owner.String(),
				Inline: false,
			},
			{
				Name:   "Source code",
				Value:  "[GitHub](https://github.com/Starshine113/berry)\n/ Licensed under the [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html)",
				Inline: false,
			},
		},
	}

	_, err = ctx.Send(embed)
	return
}

func invite(ctx *crouter.Ctx) string {
	// perms is the list of permissions the bot will be granted by default
	var perms = discordgo.PermissionReadMessages +
		discordgo.PermissionReadMessageHistory +
		discordgo.PermissionSendMessages +
		discordgo.PermissionManageMessages +
		discordgo.PermissionEmbedLinks +
		discordgo.PermissionAttachFiles +
		discordgo.PermissionUseExternalEmojis +
		discordgo.PermissionAddReactions

	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&scope=bot", ctx.Session.State.User.ID, perms)
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
