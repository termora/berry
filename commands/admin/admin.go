package admin

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

// Admin ...
type Admin struct {
	*bot.Bot

	guilds []discord.Guild

	stopStatus chan bool

	WebhookClient *webhook.Client
}

// Init ...
func Init(bot *bot.Bot) (m string, out []*bcr.Command) {
	c := &Admin{Bot: bot}
	c.stopStatus = make(chan bool, 1)

	if c.Config.Bot.TermLog.ID.IsValid() {
		c.WebhookClient = webhook.New(c.Config.Bot.TermLog.ID, c.Config.Bot.TermLog.Token)
	}

	admins := bot.Router.RequireRole("Bot Admin", c.Config.Bot.Permissions.Admins...)
	directors := bot.Router.RequireRole("Director", append(c.Config.Bot.Permissions.Admins, c.Config.Bot.Permissions.Directors...)...)

	a := bot.Router.AddCommand(&bcr.Command{
		Name:    "admin",
		Summary: "Admin commands",

		CustomPermissions: directors,
		Hidden:            true,

		Command: func(ctx *bcr.Context) (err error) {
			_, err = ctx.Send("Nothing to see here, move along...\n(Hint: use `admin help` for a list of subcommands!)", nil)
			return err
		},
	})

	a.AddSubcommand(&bcr.Command{
		Name:        "addterm",
		Aliases:     []string{"add-term"},
		Summary:     "Add a term",
		Description: "Add a term. Separate names with newlines",
		Usage:       "<names...>",
		Args:        bcr.MinArgs(1),

		CustomPermissions: directors,

		Command: c.addTerm,
	}).AddSubcommand(&bcr.Command{
		Name:    "all-in-one",
		Aliases: []string{"aio", "allinone"},

		Summary: "Add a term by passing all parameters to the command invocation",
		Usage:   "<name> <category> <description> <aliases, comma separated> <source>",

		CustomPermissions: admins,
		Command:           c.aio,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "delterm",
		Aliases: []string{"del-term"},
		Summary: "Delete a term",
		Usage:   "<id>",

		CustomPermissions: admins,
		Command:           c.delTerm,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addcategory",
		Aliases: []string{"add-category"},
		Summary: "Add a category",
		Usage:   "<name>",

		CustomPermissions: admins,
		Command:           c.addCategory,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "add-pronouns",
		Aliases: []string{"addpronouns"},
		Summary: "Add a pronoun set",
		Usage:   "<subjective>/<objective>/<poss. determiner>/<poss. pronoun>/<reflexive>",

		CustomPermissions: directors,
		Command:           c.addPronouns,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addexplanation",
		Aliases: []string{"add-explanation"},
		Summary: "Add an explanation",
		Usage:   "<names...>(newline)<explanation>",

		CustomPermissions: directors,
		Command:           c.addExplanation,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "toggleexplanationcmd",
		Aliases: []string{"toggle-explanation-cmd"},
		Summary: "Set whether or not this explanation can be triggered as a command",
		Usage:   "<id> <bool>",

		CustomPermissions: admins,
		Command:           c.toggleExplanationCmd,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setflags",
		Summary: "Set a term's flags",
		Usage:   "<id> <flag mask>",

		CustomPermissions: directors,
		Command:           c.setFlags,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setcw",
		Summary: "Set a term's CW. Use `-clear` to clear.",
		Usage:   "<id> <content warning>",

		CustomPermissions: directors,
		Command:           c.setCW,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setnote",
		Summary: "Set a term's note",
		Usage:   "<id> <note>",

		CustomPermissions: directors,
		Command:           c.setNote,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "editterm",
		Aliases: []string{"edit-term"},
		Summary: "Edit a term",
		Usage:   "<part to edit> <id> <text>",

		CustomPermissions: directors,
		Command:           c.editTerm,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "update",
		Summary: "Update the bot",

		OwnerOnly: true,
		Command:   c.update,
	})

	a.AddSubcommand(&bcr.Command{
		Name:        "restart",
		Summary:     "Restart the bot",
		Description: "If the `-s`/`--silent` flag is set, don't change the bot's status",
		Usage:       "[delay]",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("silent", "s", false, "If this flag is used, don't set the bot's status")
			return fs
		},

		OwnerOnly: true,
		Command:   c.restart,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "error",
		Summary: "Get an error by ID",
		Usage:   "<error ID>",

		CustomPermissions: admins,
		Command:           c.error,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "guilds",
		Summary: "Get a list of all guilds and their owners",

		OwnerOnly: true,
		Command:   c.cmdGuilds,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "changelog",
		Summary: "Show a list of terms added since `date`.\n`date` must be in `yyyy-mm-dd` format.",
		Usage:   "<channel> <date>",

		CustomPermissions: admins,
		Command:           c.changelog,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "update-tags",
		Summary: "Bulk update a list of term's tags. Input in CSV format",

		CustomPermissions: admins,
		Command:           c.updateTags,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "import",
		Summary: "Add a term from a correctly formatted message.",
		Usage:   "<message link|ID>",
		Args:    bcr.MinArgs(1),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.StringP("category", "c", "", "Category")
			fs.BoolP("raw-source", "r", false, "Use the provided source as-is")
			return fs
		},

		CustomPermissions: directors,
		Command:           c.importFromMessage,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "upload",
		Summary: "Upload a file",

		CustomPermissions: directors,
		Command:           c.upload,
	})

	i := bot.Router.AddCommand(bot.Router.AliasMust("ai", nil, []string{"admin", "import"}, nil))
	i.Args = bcr.MinArgs(1)

	// set status
	// this is in admin because of the `guild` owner command
	var o sync.Once
	bot.Router.State.AddHandler(func(d *gateway.ReadyEvent) {
		o.Do(func() {
			c.setStatusLoop(bot.Router.State)
		})
	})

	out = append(out, a)
	return "Bot admin commands", out
}
