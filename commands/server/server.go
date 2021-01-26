package server

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/bot"
	"github.com/starshine-sys/berry/db"
)

type commands struct {
	db *db.Db
}

// Init ...
func Init(bot *bot.Bot) (m string, out []*bcr.Command) {
	c := &commands{db: bot.DB}

	g := bot.Router.AddCommand(&bcr.Command{
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

	return "Server management commands", append(out, g)
}
