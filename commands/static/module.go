package static

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
	"github.com/termora/berry/bot/cc"
)

// Commands ...
type Bot struct {
	*bot.Bot

	start time.Time

	memberMu             sync.RWMutex
	SupportServerMembers map[discord.UserID]discord.Member
}

// Init ...
func Init(b *bot.Bot) (m string, o []*bcr.Command) {
	bot := &Bot{
		Bot:                  b,
		start:                time.Now().UTC(),
		SupportServerMembers: map[discord.UserID]discord.Member{},
	}

	bot.Router.AddHandler(bot.interactionCreate)

	submit := &bcr.Group{
		Name:        "submit",
		Description: "Submit something to the devs!",
		Subcommands: []*bcr.Command{
			{
				Name:          "feedback",
				Summary:       "Send feedback to the developers!",
				Blacklistable: false,
				SlashCommand:  bot.submitFeedback,
			},
			// {
			// 	Name:          "term",
			// 	Summary:       "Submit a term!",
			// 	Blacklistable: false,
			// 	SlashCommand:  bot.submitFeedback,
			// },
			// {
			// 	Name:          "pronouns",
			// 	Summary:       "Submit a pronoun set!",
			// 	Cooldown:      1 * time.Second,
			// 	Blacklistable: true,
			// 	SlashCommand:  bot.submitFeedback,
			// },
		},
	}

	bot.Router.AddGroup(submit)

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "ping",

		Summary:  "Check the bot's message latency",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  bot.ping,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name:    "stats",
		Aliases: []string{"about"},

		Summary:  "Show some statistics about the bot!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  bot.about,
		Options:       &[]discord.CommandOption{},
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "feedback",

		Summary:  "Send feedback to the developers!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       bot.feedback,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "credits",

		Summary:  "A list of people who helped create the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       bot.credits,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name:    "hello",
		Aliases: []string{"Hi"},

		Summary:  "Say hi!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       bot.hello,
	}))

	help := bot.Router.AddCommand(&bcr.Command{
		Name:    "help",
		Aliases: []string{"h"},

		Summary:     "Show info about how to use the bot",
		Description: "Show info about how to use the bot. If a command name is given as an argument, show the help for that command.",
		Usage:       "[command]",
		Cooldown:    1 * time.Second,

		Blacklistable: true,
		SlashCommand:  bot.help,
		Options:       &[]discord.CommandOption{},
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "permissions",
		Aliases: []string{"perms"},

		Summary:  "Show an explanation of the permissions the bot needs",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  bot.perms,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "privacy",
		Aliases: []string{"privacy-policy"},

		Summary:  "Show an explanation of the data the bot collects (and what it doesn't)",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  bot.privacy,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "commands",
		Aliases: []string{"cmds"},

		Summary:  "Show a list of all commands",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       bot.commandList,
	})

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "invite",

		Summary:  "Get an invite link",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  bot.cmdInvite,
		Options:       &[]discord.CommandOption{},
	}))

	export := bot.Router.AddCommand(&bcr.Command{
		Name:    "export",
		Summary: "Export all terms in a DM",
		Usage:   "[--compress|-x]",

		Command: bot.export,
	})

	export.AddSubcommand(&bcr.Command{
		Name:    "csv",
		Summary: "Export terms as a CSV file",

		Command: bot.exportCSV,
	})

	export.AddSubcommand(&bcr.Command{
		Name:    "xlsx",
		Summary: "Export terms as a XLSX file",

		Command: bot.exportXLSX,
	})

	o = append(o, help, export)

	// thing
	if _, err := os.Stat("custom_commands.json"); err != nil {
		return "Bot info commands", o
	}

	bytes, err := os.ReadFile("custom_commands.json")
	if err != nil {
		log.Fatalf("Error reading custom commands file: %v", err)
	}

	cmds, err := cc.ParseBytes(bytes)
	if err != nil {
		log.Fatalf("Error parsing custom commands file: %v", err)
	}

	for _, c := range cmds {
		o = append(o, bot.Router.AddCommand(c))
	}

	return "Bot info commands", o
}
