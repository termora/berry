package static

import (
	"os"
	"time"

	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
	"github.com/termora/berry/bot/cc"
)

// Commands ...
type Commands struct {
	*bot.Bot

	start time.Time
}

// Init ...
func Init(bot *bot.Bot) (m string, o []*bcr.Command) {
	c := &Commands{Bot: bot, start: time.Now()}
	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "ping",

		Summary:  "Check the bot's message latency",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.ping,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "about",

		Summary:  "Some info about the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.about,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "credits",

		Summary:  "A list of people who helped create the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.credits,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name:    "hello",
		Aliases: []string{"Hi"},

		Summary:  "Say hi!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.hello,
	}))

	help := bot.Router.AddCommand(&bcr.Command{
		Name:    "help",
		Aliases: []string{"h"},

		Summary:     "Show info about how to use the bot",
		Description: "Show info about how to use the bot. If a command name is given as an argument, show the help for that command.",
		Usage:       "[command]",
		Cooldown:    1 * time.Second,

		Blacklistable: true,
		Command:       c.help,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "permissions",
		Aliases: []string{"perms"},

		Summary:  "Show an explanation of the permissions the bot needs",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.perms,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "privacy",
		Aliases: []string{"privacy-policy"},

		Summary:  "Show an explanation of the data the bot collects (and what it doesn't)",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.privacy,
	})

	help.AddSubcommand(&bcr.Command{
		Name: "autopost",

		Summary:  "Show how to set up automatically posting terms in a channel",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.autopost,
	})

	help.AddSubcommand(&bcr.Command{
		Name: "commands",

		Summary:  "Show a list of all commands",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       bot.CommandList,
	})

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "invite",

		Summary:  "Get an invite link",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.cmdInvite,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name:    "export",
		Summary: "Export all terms in a DM",
		Usage:   "[--compress|-x]",

		Cooldown: time.Minute,

		Command: c.export,
	}))

	o = append(o, help)

	// thing
	if _, err := os.Stat("custom_commands.json"); err != nil {
		return "Static commands", o
	}

	bytes, err := os.ReadFile("custom_commands.json")
	if err != nil {
		bot.Sugar.Fatalf("Error reading custom commands file: %v", err)
	}

	cmds, err := cc.ParseBytes(bytes)
	if err != nil {
		bot.Sugar.Fatalf("Error parsing custom commands file: %v", err)
	}

	for _, c := range cmds {
		o = append(o, bot.Router.AddCommand(c))
	}

	return "Bot info commands", o
}
