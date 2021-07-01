package server

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

type commands struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (m string, out []*bcr.Command) {
	c := &commands{Bot: bot}

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

	prefixes := bot.Router.AddCommand(&bcr.Command{
		Name:    "prefixes",
		Aliases: []string{"prefix"},
		Summary: "Show this server's prefixes",

		Blacklistable: true,
		Command:       c.prefixes,
	})

	prefixes.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add the given prefix",
		Usage:   "<prefix>",

		Permissions: discord.PermissionManageGuild,
		Command:     c.addPrefix,
	})

	prefixes.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove the given prefix",
		Usage:   "<prefix>",

		Permissions: discord.PermissionManageGuild,
		Command:     c.removePrefix,
	})

	return "Server configuration commands", append(out, g, prefixes)
}
