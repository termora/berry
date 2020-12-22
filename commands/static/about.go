package static

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/termbot/db"
	"github.com/bwmarrin/discordgo"
)

var botVersion = "v0.1"

func (c *commands) about(ctx *crouter.Ctx) (err error) {
	owner, err := ctx.Session.User(c.config.Bot.BotOwners[0])
	if err != nil {
		return ctx.CommandError(err)
	}

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
				Name:   "Author",
				Value:  owner.Mention() + " / " + owner.String(),
				Inline: false,
			},
			{
				Name:   "Source code",
				Value:  "[GitHub](https://github.com/Starshine113/covebotnt) / Licensed under the [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html)",
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
		discordgo.PermissionUseExternalEmojis +
		discordgo.PermissionAddReactions

	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&scope=bot", ctx.Session.State.User.ID, perms)
}
