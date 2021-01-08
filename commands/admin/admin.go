package admin

import (
	"time"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/structs"
	"github.com/diamondburned/arikawa/v2/discord"
	"go.uber.org/zap"
)

type commands struct {
	db     *db.Db
	config *structs.BotConfig
	sugar  *zap.SugaredLogger

	admins []string
}

// Init ...
func Init(db *db.Db, sugar *zap.SugaredLogger, conf *structs.BotConfig, r *bcr.Router) {
	c := &commands{db: db, config: conf, sugar: sugar}

	r.AddCommand(&bcr.Command{
		Name:    "addterm",
		Aliases: []string{"add-term"},
		Summary: "Add a term",

		CustomPermissions: c.checkOwner,

		Command: c.addTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "delterm",
		Aliases: []string{"del-term"},
		Summary: "Delete a term",

		CustomPermissions: c.checkOwner,

		Command: c.delTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "addcategory",
		Aliases: []string{"add-category"},
		Summary: "Add a category",
		Usage:   "<name>",

		CustomPermissions: c.checkOwner,

		Command: c.addCategory,
	})

	r.AddCommand(&bcr.Command{
		Name:    "addexplanation",
		Aliases: []string{"add-explanation"},
		Summary: "Add an explanation",

		CustomPermissions: c.checkOwner,

		Command: c.addExplanation,
	})

	r.AddCommand(&bcr.Command{
		Name:    "setflags",
		Summary: "Set a term's flags",

		CustomPermissions: c.checkOwner,

		Command: c.setFlags,
	})

	r.AddCommand(&bcr.Command{
		Name:    "setcw",
		Summary: "Set a term's CW",

		CustomPermissions: c.checkOwner,

		Command: c.setCW,
	})

	r.AddCommand(&bcr.Command{
		Name:    "setnote",
		Summary: "Set a term's note",

		CustomPermissions: c.checkOwner,

		Command: c.setNote,
	})

	r.AddCommand(&bcr.Command{
		Name:    "editterm",
		Aliases: []string{"edit-term"},
		Summary: "Edit a term",

		CustomPermissions: c.checkOwner,

		Command: c.editTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "export",
		Summary: "Export all terms",
		Usage:   "[-gz] [-channel <ChannelID/Mention>]",

		Cooldown:    time.Minute,
		Permissions: discord.PermissionManageMessages,

		Command: c.export,
	})

	r.AddCommand(&bcr.Command{
		Name:    "addadmin",
		Aliases: []string{"add-admin"},
		Summary: "Add an admin",
		Usage:   "<user ID/mention>",

		OwnerOnly: true,
		Command:   c.addAdmin,
	})

	r.AddCommand(&bcr.Command{
		Name:    "update",
		Summary: "Update the bot",

		OwnerOnly: true,
		Command:   c.update,
	})

	r.AddCommand(&bcr.Command{
		Name:    "restart",
		Summary: "Restart the bot",

		OwnerOnly: true,
		Command:   c.restart,
	})
}

func (c *commands) checkOwner(ctx *bcr.Context) (string, bool) {
	if c.config.Bot.AdminServer != "" {
		if ctx.Message.GuildID.String() != c.config.Bot.AdminServer {
			return "Bot Admin", false
		}
	}
	for _, id := range c.config.Bot.BotOwners {
		if id == ctx.Author.ID.String() {
			return "", true
		}
	}

	if len(c.admins) == 0 {
		admins, err := c.db.GetAdmins()
		if err != nil {
			c.sugar.Error("Error getting admins:", err)
			return "Bot Admin", false
		}
		c.admins = admins
	}

	for _, id := range c.admins {
		if id == ctx.Author.ID.String() {
			return "", true
		}
	}

	return "Bot Admin", false
}
