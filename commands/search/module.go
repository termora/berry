package search

import (
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

type Bot struct {
	*bot.Bot
}

// Init ...
func Init(b *bot.Bot) (m string, list []*bcr.Command) {
	bot := &Bot{Bot: b}

	// add autocomplete handler
	bot.Router.AddHandler(bot.doAutocomplete)

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "search",
		Aliases: []string{"s"},

		Summary:     "Search for a term",
		Description: "Search for a term. Prefix your search with `!` to show the first result.\nUse the `-c` flag to limit search results to a specific category, and use `-i` to ignore specific tags. Use `-no-cw` to hide all terms with a CW.",
		Usage:       "[-c <category>] [-i tags] [-no-cw] <search term>",

		Blacklistable: true,

		Command: bot.search,

		SlashCommand: bot.searchSlash,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:   "query",
				Description:  "The term to search for",
				Required:     true,
				Autocomplete: true,
			},
			&discord.StringOption{
				OptionName:  "category",
				Description: "The category to limit your search to",
				Required:    false,
			},
			&discord.StringOption{
				OptionName:  "ignore-tags",
				Description: "Tags to ignore (comma-separated)",
				Required:    false,
			},
			&discord.BooleanOption{
				OptionName:  "no-cw",
				Description: "Whether to hide terms with content warnings",
				Required:    false,
			},
		},
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "random",
		Aliases: []string{"r"},

		Summary: "Show a random term (optionally filtering by category)",
		Usage:   "[category]",

		Cooldown:      time.Second,
		Blacklistable: true,

		SlashCommand: bot.random,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:  "category",
				Description: "The category to find a random term in",
			},
			&discord.StringOption{
				OptionName:  "ignore",
				Description: "The tags to ignore (comma-separated)",
			},
		},
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "explain",
		Aliases: []string{"e", "ex"},

		Summary: "Explain a topic",
		Usage:   "[explanation]",

		Cooldown:      time.Second,
		Blacklistable: false,

		SlashCommand: bot.explanation,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:   "explanation",
				Description:  "Which explanation to show",
				Required:     true,
				Autocomplete: true,
			},
		},
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:        "list",
		Summary:     "List all terms, optionally filtering by a category",
		Description: "List all terms, optionally filtering by category. Use `--full` to show a list with every term's description.",
		Usage:       "[category]",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("full", "f", false, "Show all terms' full descriptions")
			fs.BoolP("file", "F", false, "Send the list of terms as a file")
			return fs
		},

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       bot.list,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "define",
		Aliases: []string{"term", "post", "d"},
		Summary: "Post a term's definition",
		Usage:   "<term ID/name>",

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       bot.term,
		SlashCommand:  bot.termSlash,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:   "query",
				Description:  "The term to define",
				Required:     true,
				Autocomplete: true,
			},
		},
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "tag",
		Aliases: []string{"tags"},
		Summary: "Show all terms with the given tag (case-insensitive)",
		Usage:   "[tag]",

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       bot.tags,
		SlashCommand:  bot.tagsSlash,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:  "tag",
				Description: "The tag to show terms for",
				Required:    false,
			},
		},
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "files",
		Summary: "Show a list of all files in the database.",
		Usage:   "[filter]",

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       bot.files,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "file",
		Summary: "Show a file by name or ID.",
		Usage:   "<ID|name>",
		Args:    bcr.MinArgs(1),

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       bot.file,
	}))

	ap := bot.Router.AddCommand(&bcr.Command{
		Name:             "autopost",
		Summary:          "Configure the bot automatically posting terms in a channel",
		Usage:            "<channel> <interval|reset>",
		Args:             bcr.MinArgs(2),
		GuildOnly:        true,
		GuildPermissions: discord.PermissionManageGuild,

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.StringP("category", "c", "", "The category to post terms from")
			fs.StringP("role", "r", "", "The role to mention when posting a term")
			return fs
		},

		Command:      bot.autopostText,
		SlashCommand: bot.autopost,
		Options: &[]discord.CommandOption{
			discord.NewChannelOption("channel", "The channel to post to", true),
			discord.NewStringOption("interval", `How often to post a term ("reset" or "off" to disable posting in the channel)`, true),
			discord.NewStringOption("category", "The category to post terms from", false),
			discord.NewRoleOption("role", "The role to mention when posting a term", false),
		},
	})

	ap.AddSubcommand(&bcr.Command{
		Name:             "list",
		Summary:          "List this server's current autopost configuration",
		GuildPermissions: discord.PermissionManageGuild,
		GuildOnly:        true,
		Command:          bot.autopostList,
	})

	// list commands
	bot.Router.AddGroup(&bcr.Group{
		Name:        "list",
		Description: "List terms and pronouns!",
		Subcommands: []*bcr.Command{
			{
				Name:          "terms",
				Summary:       "List all terms, optionally filtering by a category",
				Blacklistable: true,
				SlashCommand:  bot.listTermsSlash,
				Options: &[]discord.CommandOption{
					discord.NewStringOption("category", "The category to list terms from", false),
					discord.NewBooleanOption("full", "Show all terms with their descriptions", false),
					discord.NewBooleanOption("file", "Send the list as a file", false),
				},
			},
			{
				Name:          "pronouns",
				Summary:       "List all pronouns",
				Blacklistable: true,
				SlashCommand:  bot.listPronounsSlash,
				Options: &[]discord.CommandOption{
					&discord.StringOption{
						OptionName:  "sort-by",
						Description: "How to sort the list",
						Required:    false,
						Choices: []discord.StringChoice{
							{Name: "Sort randomly", Value: "random"},
							{Name: "Sort alphabetically", Value: "alphabetical"},
							{Name: "Sort by number of uses", Value: "uses"},
						},
					},
				},
			},
		},
	})

	state, _ := bot.Router.StateFromGuildID(0)

	var o sync.Once
	state.AddHandler(func(_ *gateway.ReadyEvent) {
		o.Do(func() {
			go bot.autopostLoop()
		})
	})

	// aliases
	ps := bot.Router.AddCommand(bot.Router.AliasMust(
		"plural", nil,
		[]string{"search"},
		bcr.DefaultArgTransformer("-c plurality", ""),
	))
	// we need to set these manually, the default description doesn't cut it
	ps.Summary = "Search for a plurality-related term"
	ps.Description = "Search for a term in the `plurality` category. Prefix your search with `!` to show the first result."
	ps.Usage = "<search term>"

	ls := bot.Router.AddCommand(bot.Router.AliasMust(
		"lgbt", []string{"lgbtq", "l", "mogai", "queer"},
		[]string{"search"},
		bcr.DefaultArgTransformer("-c lgbtq+", ""),
	))
	// same as above
	ls.Summary = "Search for a LGBTQ+-related term"
	ls.Description = "Search for a term in the `LGBTQ+` category. Prefix your search with `!` to show the first result."
	ls.Usage = "<search term>"

	list = append(list, bot.initExplanations(bot.Router)...)
	list = append(list, ps, ls, ap)
	return "Search commands", list
}
