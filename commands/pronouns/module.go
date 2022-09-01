package pronouns

import (
	"embed"
	"strings"
	"text/template"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

//go:embed examples/*
var fs embed.FS

var templates = template.Must(template.New("").Funcs(funcs()).ParseFS(fs, "examples/*"))
var tmplCount int

// initialise number of templates
func init() {
	files, err := fs.ReadDir("examples")
	if err != nil {
		panic(err)
	}
	tmplCount = len(files)
}

type Bot struct {
	*bot.Bot

	submitCooldown *ttlcache.Cache
}

// Init ...
func Init(b *bot.Bot) (m string, list []*bcr.Command) {
	bot := &Bot{
		Bot:            b,
		submitCooldown: ttlcache.NewCache(),
	}
	bot.submitCooldown.SkipTTLExtensionOnHit(true)

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "list-pronouns",
		Aliases: []string{"pronoun-list", "listpronouns", "pronounlist"},

		Summary: "Show a list of all pronouns",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("random", "r", false, "Sort pronouns randomly")
			fs.BoolP("alphabetical", "a", false, "Sort pronouns alphabetically")
			fs.BoolP("by-uses", "u", false, "Sort pronouns by number of uses")
			return fs
		},

		Blacklistable: true,
		Cooldown:      time.Second,
		Command:       bot.list,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name: "submit-pronouns",

		Summary: "Submit a pronoun set",
		Usage:   "<pronouns, forms separated with />",

		Blacklistable: true,
		Command:       bot.submit,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "random-pronouns",
		Aliases: []string{"random-pronoun"},
		Summary: "Show a random pronoun set",
		Options: &[]discord.CommandOption{},

		Blacklistable: true,
		Cooldown:      time.Second,
		Command:       bot.random,
		SlashCommand:  bot.randomSlash,
	}))

	pronouns := bot.Router.AddCommand(&bcr.Command{
		Name:    "pronouns",
		Aliases: []string{"pronoun", "neopronoun", "neopronouns"},

		Summary: "Show pronouns (with optional name) used in a sentence",
		Usage:   "<pronouns> [name]",

		Blacklistable: true,
		Cooldown:      time.Second,
		SlashCommand:  bot.use,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:  "pronouns",
				Description: "The pronouns to show",
				Required:    true,
			},
			&discord.StringOption{
				OptionName:  "name",
				Description: "The name to use",
				Required:    false,
			},
		},
	})

	pronouns.AddSubcommand(&bcr.Command{
		Name:          "custom",
		Summary:       "Show custom pronouns that aren't in the bot",
		Usage:         "<pronoun set, space or slash separated>",
		Blacklistable: true,
		Cooldown:      time.Second,
		SlashCommand:  bot.custom,
	})

	bot.Router.AddCommand(&bcr.Command{
		Name:          "custom-pronouns",
		Summary:       "Show custom pronouns that aren't in the bot",
		Usage:         "<pronoun set, space or slash separated>",
		Blacklistable: true,
		Cooldown:      time.Second,
		SlashCommand:  bot.custom,
		Options: &[]discord.CommandOption{
			&discord.StringOption{
				OptionName:  "set",
				Description: "The pronouns to show (separated by /)",
				Required:    true,
			},
		},
	})

	pronouns.AddSubcommand(bot.Router.AliasMust("list", []string{"l"}, []string{"list-pronouns"}, nil))
	pronouns.AddSubcommand(bot.Router.AliasMust("submit", nil, []string{"submit-pronouns"}, nil))
	pronouns.AddSubcommand(bot.Router.AliasMust("random", []string{"r"}, []string{"random-pronouns"}, nil))

	bot.Router.AddHandler(bot.reactionAdd)

	return "Pronoun commands", append(list, pronouns)
}

func funcs() map[string]interface{} {
	return map[string]interface{}{
		"title": strings.Title,
	}
}
