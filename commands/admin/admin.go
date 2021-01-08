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
	c := &commands{db: db, config: conf}

	r.AddCommand(&bcr.Command{
		Name:    "AddTerm",
		Summary: "Add a term",

		CustomPermissions: c.checkOwner,

		Command: c.addTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "DelTerm",
		Summary: "Delete a term",

		CustomPermissions: c.checkOwner,

		Command: c.delTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "AddCategory",
		Summary: "Add a category",
		Usage:   "<name>",

		CustomPermissions: c.checkOwner,

		Command: c.addCategory,
	})

	r.AddCommand(&bcr.Command{
		Name:    "AddExplanation",
		Summary: "Add an explanation",

		CustomPermissions: c.checkOwner,

		Command: c.addExplanation,
	})

	r.AddCommand(&bcr.Command{
		Name:    "SetFlags",
		Summary: "Set a term's flags",

		CustomPermissions: c.checkOwner,

		Command: c.setFlags,
	})

	r.AddCommand(&bcr.Command{
		Name:    "SetCW",
		Summary: "Set a term's CW",

		CustomPermissions: c.checkOwner,

		Command: c.setCW,
	})

	r.AddCommand(&bcr.Command{
		Name:    "EditTerm",
		Summary: "Edit a term",

		CustomPermissions: c.checkOwner,

		Command: c.editTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:    "Export",
		Summary: "Export all terms",
		Usage:   "[-gz] [-channel <ChannelID/Mention>]",

		Cooldown:    time.Minute,
		Permissions: discord.PermissionManageMessages,

		Command: c.export,
	})

	r.AddCommand(&bcr.Command{
		Name:    "AddAdmin",
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
		Command:   c.update,
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
