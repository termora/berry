package server

import (
	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
)

type commands struct {
	db *db.Db
}

// Init ...
func Init(db *db.Db, r *bcr.Router) {
	c := &commands{db: db}

	g := r.AddCommand(&bcr.Command{
		Name:    "blacklist",
		Aliases: []string{"bl", "blocklist", "blocking"},
		Summary: "Manage this server's blacklist",

		Permissions: discord.PermissionManageGuild,
		Command:     c.blacklist,
	})

	g.AddSubcommand(&bcr.Command{
		Name:    "add",
		Aliases: []string{"block"},
		Summary: "Add channels to the blacklist",
		Usage:   "<channels...>",

		Permissions: discord.PermissionManageGuild,
		Command:     c.blacklistAdd,
	})

	g.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Aliases: []string{"rm", "unblock"},
		Summary: "Remove a channel from the blacklist",
		Usage:   "<channel>",

		Permissions: discord.PermissionManageGuild,
		Command:     c.blacklistRemove,
	})
}
