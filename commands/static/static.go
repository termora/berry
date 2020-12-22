package static

import (
	"github.com/Starshine113/crouter"
	"github.com/Starshine113/termbot/structs"
)

type commands struct {
	config *structs.BotConfig
}

// Init ...
func Init(conf *structs.BotConfig, r *crouter.Router) {
	c := &commands{config: conf}
	r.AddCommand(&crouter.Command{
		Name: "Ping",

		Summary: "Check the bot's message latency",

		Blacklistable: true,
		Command:       c.ping,
	})

	r.AddCommand(&crouter.Command{
		Name: "About",

		Summary: "Some info about the bot",

		Blacklistable: true,
		Command:       c.about,
	})

	r.AddCommand(&crouter.Command{
		Name:    "Hello",
		Aliases: []string{"Hi"},

		Summary: "Some hi!",

		Blacklistable: true,
		Command:       c.hello,
	})

	r.AddCommand(&crouter.Command{
		Name: "Help",

		Summary: "Show info about how to use the bot",

		Blacklistable: true,
		Command:       c.help,
	})

	r.AddCommand(&crouter.Command{
		Name: "Invite",

		Summary: "Get an invite link",

		Blacklistable: true,
		Command:       c.cmdInvite,
	})
}
