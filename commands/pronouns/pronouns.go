package pronouns

import (
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

var templates = template.Must(template.New("").Funcs(funcs()).ParseGlob("pronoun-examples/*"))
var tmplCount int

// initialise number of templates
func init() {
	files, err := ioutil.ReadDir("pronoun-examples")
	if err != nil {
		panic(err)
	}
	tmplCount = len(files)
}

type commands struct {
	*bot.Bot

	submitCooldown *ttlcache.Cache
}

// Init ...
func Init(bot *bot.Bot) (m string, list []*bcr.Command) {
	c := &commands{
		Bot:            bot,
		submitCooldown: ttlcache.NewCache(),
	}
	c.submitCooldown.SkipTTLExtensionOnHit(true)

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "list-pronouns",
		Aliases: []string{"pronoun-list", "listpronouns", "pronounlist"},

		Summary: "Show a list of all pronouns",

		Blacklistable: true,
		Cooldown:      time.Second,
		Command:       c.list,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name: "submit-pronouns",

		Summary: "Submit a pronoun set",
		Usage:   "<pronouns, forms separated with />",

		Blacklistable: true,
		Command:       c.submit,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "random-pronouns",
		Aliases: []string{"random-pronoun"},
		Summary: "Show a random pronoun set",

		Blacklistable: true,
		Cooldown:      time.Second,
		Command:       c.random,
	}))

	pronouns := bot.Router.AddCommand(&bcr.Command{
		Name:    "pronouns",
		Aliases: []string{"pronoun", "neopronoun", "neopronouns"},

		Summary: "Show pronouns (with optional name) used in a sentence",
		Usage:   "<pronouns> [name]",

		Blacklistable: true,
		Cooldown:      time.Second,
		Command:       c.use,
	})

	pronouns.AddSubcommand(&bcr.Command{
		Name:          "custom",
		Summary:       "Show custom pronouns that aren't in the bot",
		Usage:         "<pronoun set, space or slash separated>",
		Blacklistable: true,
		Cooldown:      time.Second,
		Command:       c.custom,
	})

	pronouns.AddSubcommand(bot.Router.AliasMust("list", []string{"l"}, []string{"list-pronouns"}, nil))
	pronouns.AddSubcommand(bot.Router.AliasMust("submit", nil, []string{"submit-pronouns"}, nil))
	pronouns.AddSubcommand(bot.Router.AliasMust("random", []string{"r"}, []string{"random-pronouns"}, nil))

	bot.Router.State.AddHandler(c.reactionAdd)

	return "Pronoun commands", append(list, pronouns)
}

func funcs() map[string]interface{} {
	return map[string]interface{}{
		"title": strings.Title,
	}
}
