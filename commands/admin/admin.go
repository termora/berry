package admin

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/bot"
	"github.com/starshine-sys/berry/db"
	"github.com/starshine-sys/berry/structs"
	"go.uber.org/zap"
)

// Admin ...
type Admin struct {
	db     *db.Db
	config *structs.BotConfig
	sugar  *zap.SugaredLogger

	admins []string

	guilds []discord.Guild

	stopStatus chan bool
}

func (Admin) String() string {
	return "Bot Admin"
}

// Check ...
func (c *Admin) Check(ctx *bcr.Context) (bool, error) {
	if c.config.Bot.AdminServers != nil {
		var inServer bool
		for _, s := range c.config.Bot.AdminServers {
			if ctx.Message.GuildID == s {
				inServer = true
				break
			}
		}
		if !inServer {
			return false, nil
		}
	}
	for _, id := range c.config.Bot.BotOwners {
		if id == ctx.Author.ID {
			return true, nil
		}
	}

	if len(c.admins) == 0 {
		admins, err := c.db.GetAdmins()
		if err != nil {
			c.sugar.Error("Error getting admins:", err)
			return false, err
		}
		c.admins = admins
	}

	for _, id := range c.admins {
		if id == ctx.Author.ID.String() {
			return true, nil
		}
	}

	return false, nil
}

// Init ...
func Init(bot *bot.Bot) (m string, out []*bcr.Command) {
	c := &Admin{db: bot.DB, config: bot.Config, sugar: bot.Sugar}
	c.stopStatus = make(chan bool, 1)

	a := bot.Router.AddCommand(&bcr.Command{
		Name:    "admin",
		Summary: "Admin commands",

		CustomPermissions: c,
		Hidden:            true,

		Command: func(ctx *bcr.Context) (err error) {
			_, err = ctx.Send("Nothing to see here, move along...\n(Hint: use `admin help` for a list of subcommands!)", nil)
			return err
		},
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addterm",
		Aliases: []string{"add-term"},
		Summary: "Add a term",

		CustomPermissions: c,

		Command: c.addTerm,
	}).AddSubcommand(&bcr.Command{
		Name:    "all-in-one",
		Aliases: []string{"aio", "allinone"},

		Summary: "Add a term by passing all parameters to the command invocation",
		Usage:   "<name> <category> <description> <aliases, comma separated> <source>",

		CustomPermissions: c,
		Command:           c.aio,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "delterm",
		Aliases: []string{"del-term"},
		Summary: "Delete a term",
		Usage:   "<id>",

		CustomPermissions: c,

		Command: c.delTerm,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addcategory",
		Aliases: []string{"add-category"},
		Summary: "Add a category",
		Usage:   "<name>",

		CustomPermissions: c,

		Command: c.addCategory,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "add-pronouns",
		Aliases: []string{"addpronouns"},
		Summary: "Add a pronoun set",
		Usage:   "<subjective>/<objective>/<poss. determiner>/<poss. pronoun>/<reflexive>",

		CustomPermissions: c,

		Command: c.addPronouns,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addexplanation",
		Aliases: []string{"add-explanation"},
		Summary: "Add an explanation",
		Usage:   "<names...>(newline)<explanation>",

		CustomPermissions: c,

		Command: c.addExplanation,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "toggleExplanationCmd",
		Aliases: []string{"toggle-explanation-cmd"},
		Summary: "Set whether or not this explanation can be triggered as a command",
		Usage:   "<id> <bool>",

		CustomPermissions: c,

		Command: c.toggleExplanationCmd,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setflags",
		Summary: "Set a term's flags",
		Usage:   "<id> <flag mask>",

		CustomPermissions: c,

		Command: c.setFlags,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setcw",
		Summary: "Set a term's CW",
		Usage:   "<id> <content warning>",

		CustomPermissions: c,

		Command: c.setCW,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setnote",
		Summary: "Set a term's note",
		Usage:   "<id> <note>",

		CustomPermissions: c,

		Command: c.setNote,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "editterm",
		Aliases: []string{"edit-term"},
		Summary: "Edit a term",
		Usage:   "<part to edit> <id> <text>",

		CustomPermissions: c,

		Command: c.editTerm,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addadmin",
		Aliases: []string{"add-admin"},
		Summary: "Add an admin",
		Usage:   "<user ID/mention>",

		OwnerOnly: true,
		Command:   c.addAdmin,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "update",
		Summary: "Update the bot",

		OwnerOnly: true,
		Command:   c.update,
	})

	a.AddSubcommand(&bcr.Command{
		Name:        "restart",
		Summary:     "Restart the bot",
		Description: "If the `-s`/`--silent` flag is set, don't change the bot's status",
		Usage:       "[delay] [-s/--silent]",

		OwnerOnly: true,
		Command:   c.restart,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "error",
		Summary: "Get an error by ID",
		Usage:   "<error ID>",

		CustomPermissions: c,
		Command:           c.error,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "guilds",
		Summary: "Get a list of all guilds and their owners",

		OwnerOnly: true,
		Command:   c.cmdGuilds,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "changelog",
		Summary: "Show a list of terms added since `date`.\n`date` must be in `yyyy-mm-dd` format.",
		Usage:   "<channel> <date>",

		CustomPermissions: c,
		Command:           c.changelog,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "token",
		Summary: "Get an API token",

		CustomPermissions: c,
		Command:           c.token,
	}).AddSubcommand(&bcr.Command{
		Name:    "refresh",
		Summary: "Refresh your API token",

		CustomPermissions: c,
		Command:           c.refreshToken,
	})

	// set status
	// this is in admin because of the `guild` owner command
	var o sync.Once
	bot.Router.Session.AddHandler(func(d *gateway.ReadyEvent) {
		o.Do(func() {
			c.setStatusLoop(bot.Router.Session)
		})
	})

	out = append(out, a)
	return "Bot admin commands", out
}
