package static

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Commands) perms(ctx bcr.Contexter) (err error) {
	return ctx.SendX("", discord.Embed{
		Title:       "Required permissions",
		Description: c.permText(ctx),
		Color:       db.EmbedColour,
	})
}

func (c *Commands) permText(ctx bcr.Contexter) string {
	return fmt.Sprintf(`%v requires the following permissions to function correctly:
- **Read Messages** & **Send Messages**: to respond to commands
- **Read Message History**: for the %vsearch command to work
- **Manage Messages**: to delete reactions on menus
- **Embed Links**: to send responses for most commands
- **Add Reactions**: for menus to work
- **Use External Emojis**: to use custom emotes in a couple of commands`, c.Router.Bot.Username, c.Config.Bot.Prefixes[0])
}

func (c *Commands) privacy(ctx bcr.Contexter) (err error) {
	return ctx.SendX("", discord.Embed{
		Title:       "Privacy",
		Description: c.privacyText(ctx),
		Color:       db.EmbedColour,
	})
}

func (c *Commands) privacyText(ctx bcr.Contexter) string {
	return fmt.Sprintf(`We're not lawyers, we don't want to write a document that no one will (or even can) read.

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
	
	%v is open source, and its source code is available [on GitHub](%v). While we cannot *prove* that this is the code powering the bot, we promise that it is.`, c.Router.Bot.Username, c.Router.Bot.Username, c.Router.Bot.Username, c.Router.Bot.Username, c.Router.Bot.Username, c.Config.Bot.Website, c.Router.Bot.Username, c.Config.Bot.Git)
}

func (c *Commands) autopost(ctx bcr.Contexter) (err error) {
	_, err = ctx.Send("", discord.Embed{
		Title:       "Autopost",
		Description: c.autopostText(ctx),
		Color:       db.EmbedColour,
	})
	return
}

func (c *Commands) autopostText(ctx bcr.Contexter) string {
	return fmt.Sprintf("To automatically post terms at a set interval, you can use the `/autopost` (or `%vautopost`) command. Check out `%vautopost help` for how to use it.\n\nNote: this command previously recommended using a bot such as YAGPDB.xyz to automatically post terms. This still works for now, but please switch over to the built-in command.", c.Config.Bot.Prefixes[0], c.Config.Bot.Prefixes[0])
}

func (c *Commands) help(ctx bcr.Contexter) (err error) {
	// help for commands
	if v, ok := ctx.(*bcr.Context); ok {
		if len(v.Args) > 0 {
			return v.Help(v.Args)
		}
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
				Value: fmt.Sprintf("`pronouns`: see how pronouns are used in a sentence! (optionally with your name)\n`pronouns list`: list all pronouns known to %v!\n`pronouns submit`: submit a pronoun set to be added!", c.Router.Bot.Username),
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

	components := discord.Components(
		&discord.SelectComponent{
			CustomID:    "help_options",
			Placeholder: "More info about...",
			Options: []discord.SelectOption{
				{
					Label:       "Permissions",
					Description: fmt.Sprintf("Show all the permissions %v needs!", c.Router.Bot.Username),
					Value:       "permissions",
				},
				{
					Label:       "Privacy",
					Description: fmt.Sprintf("Show %v's privacy policy!", c.Router.Bot.Username),
					Value:       "privacy",
				},
				{
					Label:       "Automatically posting terms",
					Description: fmt.Sprintf("How to make %v automatically post terms!", c.Router.Bot.Username),
					Value:       "autopost",
				},
			},
		},
	)

	msg, err := ctx.SendComponents(components, "", e)
	if err != nil {
		return err
	}

	rm := ctx.Session().AddHandler(func(ev *gateway.InteractionCreateEvent) {
		if ev.Message == nil || ev.Data == nil {
			return
		}

		data, ok := ev.Data.(*discord.SelectInteraction)
		if !ok {
			return
		}

		if ev.Message.ID != msg.ID {
			return
		}

		s := ""
		switch data.Values[0] {
		case "permissions":
			s = c.permText(ctx)
		case "privacy":
			s = c.privacyText(ctx)
		case "autopost":
			s = c.autopostText(ctx)
		default:
			ctx.Session().RespondInteraction(ev.ID, ev.Token, api.InteractionResponse{
				Type: api.UpdateMessage,
				Data: &api.InteractionResponseData{
					Components: &components,
				},
			})
		}

		ctx.Session().RespondInteraction(ev.ID, ev.Token, api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Embeds: &[]discord.Embed{{
					Description: s,
					Color:       db.EmbedColour,
				}},
				Flags: api.EphemeralResponse,
			},
		})
		return
	})

	time.AfterFunc(5*time.Minute, func() {
		ctx.Session().EditMessageComplex(msg.ChannelID, msg.ID, api.EditMessageData{
			Components: discord.ComponentsPtr(),
		})
		rm()
	})

	return err
}

func (c *Commands) cmdInvite(ctx bcr.Contexter) (err error) {
	s := fmt.Sprintf("Use this link to invite me to your server: <%v>", c.invite())
	if _, ok := ctx.(*bcr.Context); ok {
		s += "\n\nYou can use the `help permissions` command to get a detailed explanation of all permissions required."
	}

	return ctx.SendEphemeral(s)
}
