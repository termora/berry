package server

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

type Bot struct {
	*bot.Bot
}

// Init ...
func Init(b *bot.Bot) (m string, out []*bcr.Command) {
	bot := &Bot{Bot: b}

	g := bot.Router.AddCommand(&bcr.Command{
		Name:    "blacklist",
		Aliases: []string{"bl", "blocklist", "blocking"},
		Summary: "Manage this server's blacklist",

		Permissions: discord.PermissionManageGuild,
		Command:     bot.blacklist,
	})

	g.AddSubcommand(&bcr.Command{
		Name:    "add",
		Aliases: []string{"block"},
		Summary: "Add channels to the blacklist",
		Usage:   "<channels...>",

		Permissions: discord.PermissionManageGuild,
		Command:     bot.blacklistAdd,
	})

	g.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Aliases: []string{"rm", "unblock"},
		Summary: "Remove a channel from the blacklist",
		Usage:   "<channel>",

		Permissions: discord.PermissionManageGuild,
		Command:     bot.blacklistRemove,
	})

	prefixes := bot.Router.AddCommand(&bcr.Command{
		Name:    "prefixes",
		Aliases: []string{"prefix"},
		Summary: "Show this server's prefixes",

		Blacklistable: true,
		Command:       bot.prefixes,
	})

	prefixes.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add the given prefix",
		Usage:   "<prefix>",

		Permissions: discord.PermissionManageGuild,
		Command:     bot.addPrefix,
	})

	prefixes.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove the given prefix",
		Usage:   "<prefix>",

		Permissions: discord.PermissionManageGuild,
		Command:     bot.removePrefix,
	})

	return "Server configuration commands", append(out, g, prefixes)
}
