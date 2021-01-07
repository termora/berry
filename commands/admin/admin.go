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
		Name:        "AddTerm",
		Description: "Add a term",

		CustomPermissions: c.checkOwner,

		Command: c.addTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:        "DelTerm",
		Description: "Delete a term",

		CustomPermissions: c.checkOwner,

		Command: c.delTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:        "AddCategory",
		Description: "Add a category",
		Usage:       "<name>",

		CustomPermissions: c.checkOwner,

		Command: c.addCategory,
	})

	r.AddCommand(&bcr.Command{
		Name:        "AddExplanation",
		Description: "Add an explanation",

		CustomPermissions: c.checkOwner,

		Command: c.addExplanation,
	})

	r.AddCommand(&bcr.Command{
		Name:        "SetFlags",
		Description: "Set a term's flags",

		CustomPermissions: c.checkOwner,

		Command: c.setFlags,
	})

	r.AddCommand(&bcr.Command{
		Name:        "SetCW",
		Description: "Set a term's CW",

		CustomPermissions: c.checkOwner,

		Command: c.setCW,
	})

	r.AddCommand(&bcr.Command{
		Name:        "EditTerm",
		Description: "Edit a term",

		CustomPermissions: c.checkOwner,

		Command: c.editTerm,
	})

	r.AddCommand(&bcr.Command{
		Name:        "Export",
		Description: "Export all terms",
		Usage:       "[-gz] [-channel <ChannelID/Mention>]",

		Cooldown:    time.Minute,
		Permissions: discord.PermissionManageMessages,

		Command: c.export,
	})

	r.AddCommand(&bcr.Command{
		Name:        "AddAdmin",
		Description: "Add an admin",
		Usage:       "<user ID/mention>",

		OwnerOnly: true,
		Command:   c.addAdmin,
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
