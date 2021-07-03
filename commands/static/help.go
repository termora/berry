package static

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Commands) perms(ctx *bcr.Context) (err error) {
	_, err = ctx.Send("", discord.Embed{
		Title: "Required permissions",
		Description: fmt.Sprintf(`%v requires the following permissions to function correctly:
- **Read Messages** & **Send Messages**: to respond to commands
- **Read Message History**: for the %vsearch command to work
- **Manage Messages**: to delete reactions on menus
- **Embed Links**: to send responses for most commands
- **Add Reactions**: for menus to work
- **Use External Emojis**: to use custom emotes in a couple of commands`, ctx.Bot.Username, c.Config.Bot.Prefixes[0]),
		Color: db.EmbedColour,
	})
	return
}

func (c *Commands) privacy(ctx *bcr.Context) (err error) {
	_, err = ctx.Send("", discord.Embed{
		Title: "Privacy",
		Description: fmt.Sprintf(`We're not lawyers, we don't want to write a document that no one will (or even can) read.

By continuing to use %v's commands, you consent to us processing the data listed below in a manner compliant with the GDPR.

This is the data %v collects:

- A list of blacklisted channels per server

This is the data %v collects, and which is deleted after 30 days:

- Information about internal errors: command, user ID, and channel ID

This is the data %v stores in memory for as long as it's needed (up to 15 minutes), and which is always wiped on a bot restart:

- Message metadata *for its own messages*
- Message metadata for messages that trigger commands

This is the data %v does *not* collect:

- Any message contents
- Any user information
- Information about messages that do not trigger commands

Additionally, there are daily database backups, which only include a list of blacklisted channels (as well as all terms/explanations).

To delete server information from the database, simply have the bot leave the server, through kicking or banning it. Do note that this does *not* delete server information from database backups, only the live database (and any later backups). Contact us (on the support server, or [here](%vcontact)) if you want that deleted too, or have any other requests regarding your data. We'll comply with these within 30 days.

%v is open source, and its source code is available [on GitHub](%v). While we cannot *prove* that this is the code powering the bot, we promise that it is.`, ctx.Bot.Username, ctx.Bot.Username, ctx.Bot.Username, ctx.Bot.Username, ctx.Bot.Username, c.Config.Bot.Website, ctx.Bot.Username, c.Config.Bot.Git),
		Color: db.EmbedColour,
	})
	return err
}

func (c *Commands) autopost(ctx *bcr.Context) (err error) {
	_, err = ctx.Send("", discord.Embed{
		Title:       "Autopost",
		Description: fmt.Sprintf("To automatically post terms at a set interval, you can use the following custom command for [YAGPDB.xyz](https://yagpdb.xyz/):\n```{{/* Recommended trigger: Minute/Hourly interval */}}\n\n%vrandom\n{{deleteResponse 1}}```\nOther bots may have similar functionality; if you need a bot whitelisted for commands, feel free to ask on the support server.", ctx.Prefix),
		Color:       db.EmbedColour,
	})
	return
}

func (c *Commands) help(ctx *bcr.Context) (err error) {
	// help for commands
	if len(ctx.Args) > 0 {
		return ctx.Help(ctx.Args)
	}

	e := discord.Embed{
		Color: db.EmbedColour,
		Title: "Help",
		Fields: []discord.EmbedField{
			{
				Name:  "Bot info",
				Value: "`help`: show a list of commands, and some info about the bot! (alias: `h`)\n`help privacy`: show the bot's privacy policy!\n`help commands`: show the full list of commands!\n`credits`: show the people who contributed to the bot!\n`invite`: get an invite link!\n`feedback`: send feedback to the developers!",
			},
			{
				Name:  "Terms",
				Value: "`search`: search the database for a term! (alias: `s`)\n`random`: show a random term! (alias: `r`)\n`define`: show the term with the given name, or the closest match! (alias: `d`)",
			},
			{
				Name:  "Pronouns",
				Value: fmt.Sprintf("`pronouns`: see how pronouns are used in a sentence! (optionally with your name)\n`pronouns list`: list all pronouns known to %v!\n`pronouns submit`: submit a pronoun set to be added!", ctx.Bot.Username),
			},
			{
				Name:  "For staff",
				Value: "You can blacklist most commands, with the exception of `explain`, using the `blacklist` command.\nYou can also change the prefixes the bot uses with the `prefix` command.",
			},
		},
		Footer: &discord.EmbedFooter{
			Text: "Use `help <command>` for more information on a specific command.",
		},
	}

	// if custom help fields are defined, add those
	if len(c.Config.Bot.HelpFields) != 0 {
		e.Fields = append(e.Fields, c.Config.Bot.HelpFields...)
	}

	_, err = ctx.Send("", e)
	return err
}

func (c *Commands) cmdInvite(ctx *bcr.Context) (err error) {
	_, err = ctx.Sendf("Use this link to invite me to your server: <%v>\n\nYou can use the `help permissions` command to get a detailed explanation of all permissions required.", c.invite(ctx))
	return
}
