package static

import (
	"time"

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
		Name: "ping",

		Summary:  "Check the bot's message latency",
		Cooldown: 3 * time.Second,

		Blacklistable: true,
		Command:       c.ping,
	})

	r.AddCommand(&crouter.Command{
		Name: "about",

		Summary:  "Some info about the bot",
		Cooldown: 5 * time.Second,

		Blacklistable: true,
		Command:       c.about,
	})

	r.AddCommand(&crouter.Command{
		Name:    "hello",
		Aliases: []string{"Hi"},

		Summary:  "Say hi!",
		Cooldown: 3 * time.Second,

		Blacklistable: true,
		Command:       c.hello,
	})

	r.AddCommand(&crouter.Command{
		Name: "help",

		Summary:  "Show info about how to use the bot",
		Cooldown: 5 * time.Second,

		Blacklistable: true,
		Command:       c.help,
	})

	r.AddCommand(&crouter.Command{
		Name: "invite",

		Summary:  "Get an invite link",
		Cooldown: 5 * time.Second,

		Blacklistable: true,
		Command:       c.cmdInvite,
	})
}
