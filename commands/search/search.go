package search

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

type commands struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (m string, list []*bcr.Command) {
	c := commands{Bot: bot}

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "search",
		Aliases: []string{"s"},

		Summary:     "Search for a term",
		Description: "Search for a term. Prefix your search with `!` to show the first result.\nUse the `-c` flag to limit search results to a specific category, and use `-i` to ignore specific tags.",
		Usage:       "[-c <category>] [-h] [-i tags] <search term>",

		Blacklistable: true,

		Command: c.search,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "random",
		Aliases: []string{"r"},

		Summary: "Show a random term (optionally filtering by category)",
		Usage:   "[category]",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.StringSliceP("ignore-tags", "i", []string{}, "Specific tags (comma-separated) to ignore")
			return fs
		},

		Cooldown:      time.Second,
		Blacklistable: true,

		Command: c.random,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "explain",
		Aliases: []string{"e", "ex"},

		Summary: "Show a single explanation, or a list of all explanations",
		Usage:   "[explanation]",

		Cooldown:      time.Second,
		Blacklistable: false,

		Command: c.explanation,
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
		Command:       c.list,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "post",
		Aliases: []string{"term", "define", "d"},
		Summary: "Post a single term",
		Usage:   "<term ID/name>",

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       c.term,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "tag",
		Aliases: []string{"tags"},
		Summary: "Show all terms with the given tag (case-insensitive)",
		Usage:   "[tag]",

		Cooldown:      time.Second,
		Blacklistable: true,
		Command:       c.tags,
	}))

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
		"lgbt", []string{"lgbtq", "l"},
		[]string{"search"},
		bcr.DefaultArgTransformer("-c lgbtq+", ""),
	))
	// same as above
	ls.Summary = "Search for a LGBTQ+-related term"
	ls.Description = "Search for a term in the `LGBTQ+` category. Prefix your search with `!` to show the first result."
	ls.Usage = "<search term>"

	list = append(list, c.initExplanations(bot.Router)...)
	list = append(list, ps, ls)
	return "Search commands", list
}
