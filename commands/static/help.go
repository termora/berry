package static

import (
	"fmt"
	"strings"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/termbot/db"
	"github.com/bwmarrin/discordgo"
)

func (c *commands) help(ctx *crouter.Ctx) (err error) {
	if ctx.RawArgs == "permissions" || ctx.RawArgs == "perms" {
		_, err = ctx.Embed("Required permissions", fmt.Sprintf(`%v requires the following permissions to function correctly:
		- **Read Messages** & **Send Messages**: to respond to commands
		- **Read Message History**: for the %vsearch command to work
		- **Manage Messages**: to delete reactions on menus
		- **Embed Links**: to send responses for most commands
		- **Add Reactions**: for menus to work
		- **Use External Emojis**: to use custom emotes in a couple of commands`, ctx.BotUser.Username, ctx.Router.Prefixes[0]), db.EmbedColour)
		return
	}

	e := &discordgo.MessageEmbed{
		Color: db.EmbedColour,
		Title: "Help",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Prefixes",
				Value: fmt.Sprintf("%v uses the prefixes %v, and %v.", ctx.BotUser.Username, strings.Join(ctx.Router.Prefixes[:len(ctx.Router.Prefixes)-2], ", "), ctx.BotUser.Mention()),
			},
			{
				Name:  "Bot info",
				Value: "`help`: show a list of commands, and some info about the bot\n`about`: show more in-depth info about the bot.\n`ping`: check the bot's latency\n`hello`: say hi to the bot!\n`invite`: get an invite link for the bot",
			},
			{
				Name:  "Terms",
				Value: "`search`: search the database for a term (alias: `s`)\n`random`: show a random term (alias: `r`)\n`term`: get a term by its ID, useful for mods posting in a channel",
			},
			{
				Name:  "Explanations",
				Value: "`explain`: get a list of all registered explanations (aliases: `e`, `ex`)\n`explain <topic>`: explain the given topic",
			},
			{
				Name:  "For staff",
				Value: fmt.Sprintf("You can blacklist most commands, with the exception of `explain`, using the following commands:\n`blacklist`: show the current blacklist\n`blacklist add`: add a channel to the blacklist\n`blacklist remove`: remove a channel from the blacklist\n\nTo stop %v from responding in a channel completely, deny it the \"Read Messages\" permission in that channel.", ctx.BotUser.Username),
			},
		},
	}
	if c.config.Bot.ServerInvite != "" {
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:  "Support server",
			Value: fmt.Sprintf("Use this link to join the support server: %v", c.config.Bot.ServerInvite),
		})
	}
	_, err = ctx.Send(e)
	return err
}

func (c *commands) cmdInvite(ctx *crouter.Ctx) (err error) {
	_, err = ctx.Sendf("Use this link to invite me to your server: <%v>\n\nYou can use the `%vhelp permissions` command to get a detailed explanation of all permissions required.", invite(ctx), ctx.Router.Prefixes[0])
	return
}
