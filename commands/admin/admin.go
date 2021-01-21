package admin

import (
	"sync"
	"time"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/structs"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"go.uber.org/zap"
)

type commands struct {
	db     *db.Db
	config *structs.BotConfig
	sugar  *zap.SugaredLogger

	admins []string

	guilds []discord.Guild
}

func (commands) String() string {
	return "Bot Admin"
}

func (c *commands) Check(ctx *bcr.Context) (bool, error) {
	if c.config.Bot.AdminServer != "" {
		if ctx.Message.GuildID.String() != c.config.Bot.AdminServer {
			return false, nil
		}
	}
	for _, id := range c.config.Bot.BotOwners {
		if id == ctx.Author.ID.String() {
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
func Init(db *db.Db, sugar *zap.SugaredLogger, conf *structs.BotConfig, r *bcr.Router) {
	c := &commands{db: db, config: conf, sugar: sugar}

	a := r.AddCommand(&bcr.Command{
		Name:    "admin",
		Summary: "Admin commands",

		CustomPermissions: c,

		Command: func(ctx *bcr.Context) (err error) {
			_, err = ctx.Send("Nothing to see here, move along...", nil)
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

		CustomPermissions: c,

		Command: c.setFlags,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setcw",
		Summary: "Set a term's CW",

		CustomPermissions: c,

		Command: c.setCW,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setnote",
		Summary: "Set a term's note",

		CustomPermissions: c,

		Command: c.setNote,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "editterm",
		Aliases: []string{"edit-term"},
		Summary: "Edit a term",

		CustomPermissions: c,

		Command: c.editTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "export",
		Summary: "Export all terms",
		Usage:   "[-gz] [-channel <ChannelID/Mention>]",

		Cooldown: time.Minute,

		Command: c.export,
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
		Name:    "restart",
		Summary: "Restart the bot",

		OwnerOnly: true,
		Command:   c.restart,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "error",
		Summary: "Get an error by ID",

		CustomPermissions: c,
		Command:           c.error,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "guilds",
		Summary: "Get a list of all guilds and their owners",

		OwnerOnly: true,
		Command:   c.cmdGuilds,
	})

	token := a.AddSubcommand(&bcr.Command{
		Name:    "token",
		Summary: "Get an API token",

		CustomPermissions: c,
		Command:           c.token,
	})

	token.AddSubcommand(&bcr.Command{
		Name:    "refresh",
		Summary: "Refresh your API token",

		CustomPermissions: c,
		Command:           c.refreshToken,
	})

	// set status
	// this is in admin because of the `guild` owner command
	var o sync.Once
	r.Session.AddHandler(func(d *gateway.ReadyEvent) {
		o.Do(func() {
			c.setStatusLoop(r.Session)
		})
	})
}
