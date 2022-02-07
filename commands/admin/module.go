package admin

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
	"github.com/termora/berry/commands/admin/auditlog"
)

// Bot ...
type Bot struct {
	*bot.Bot

	stopStatus chan bool

	WebhookClient *webhook.Client

	AuditLog *auditlog.AuditLog
}

func (bot *Bot) guildCreate(ev *gateway.GuildCreateEvent) {
	bot.GuildsMu.Lock()
	bot.Guilds[ev.ID] = ev.Guild
	bot.GuildsMu.Unlock()
	return
}

func (bot *Bot) guildDelete(ev *gateway.GuildDeleteEvent) {
	bot.GuildsMu.Lock()
	delete(bot.Guilds, ev.ID)
	bot.GuildsMu.Unlock()
}

// Init ...
func Init(b *bot.Bot) (m string, out []*bcr.Command) {
	bot := &Bot{Bot: b}
	bot.stopStatus = make(chan bool, 1)
	bot.AuditLog = auditlog.New(b)

	admins := bot.Router.RequireRole("Bot Admin", bot.Config.Bot.Admins...)
	directors := bot.Router.RequireRole("Director", append(bot.Config.Bot.Admins, bot.Config.Bot.Directors...)...)

	a := bot.Router.AddCommand(&bcr.Command{
		Name:    "admin",
		Summary: "Admin commands",

		CustomPermissions: directors,
		Hidden:            true,

		Command: func(ctx *bcr.Context) (err error) {
			_, err = ctx.Send("Nothing to see here, move along...\n(Hint: use `admin help` for a list of subcommands!)")
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

		Command: bot.addTerm,
	}).AddSubcommand(&bcr.Command{
		Name:    "all-in-one",
		Aliases: []string{"aio", "allinone"},

		Summary: "Add a term by passing all parameters to the command invocation",
		Usage:   "<name> <category> <description> <aliases, comma separated> <source>",

		CustomPermissions: admins,
		Command:           bot.aio,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "delterm",
		Aliases: []string{"del-term"},
		Summary: "Delete a term",
		Usage:   "<id>",

		CustomPermissions: admins,
		Command:           bot.delTerm,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addcategory",
		Aliases: []string{"add-category"},
		Summary: "Add a category",
		Usage:   "<name>",

		CustomPermissions: admins,
		Command:           bot.addCategory,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "add-pronouns",
		Aliases: []string{"addpronouns"},
		Summary: "Add a pronoun set",
		Usage:   "<subjective>/<objective>/<poss. determiner>/<poss. pronoun>/<reflexive>",

		CustomPermissions: directors,
		Command:           bot.addPronouns,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "addexplanation",
		Aliases: []string{"add-explanation"},
		Summary: "Add an explanation",
		Usage:   "<names...>(newline)<explanation>",

		CustomPermissions: directors,
		Command:           bot.addExplanation,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "toggleexplanationcmd",
		Aliases: []string{"toggle-explanation-cmd"},
		Summary: "Set whether or not this explanation can be triggered as a command",
		Usage:   "<id> <bool>",

		CustomPermissions: admins,
		Command:           bot.toggleExplanationCmd,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setflags",
		Summary: "Set a term's flags",
		Usage:   "<id> <flag mask>",

		CustomPermissions: directors,
		Command:           bot.setFlags,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setcw",
		Summary: "Set a term's CW. Use `-clear` to clear.",
		Usage:   "<id> <content warning>",

		CustomPermissions: directors,
		Command:           bot.setCW,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "setnote",
		Summary: "Set a term's note",
		Usage:   "<id> <note>",

		CustomPermissions: directors,
		Command:           bot.setNote,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "editterm",
		Aliases: []string{"edit-term"},
		Summary: "Edit a term",
		Usage:   "<part to edit> <id> <text>",

		CustomPermissions: directors,
		Command:           bot.editTerm,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "update",
		Summary: "Update the bot",

		OwnerOnly: true,
		Command:   bot.update,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "error",
		Summary: "Get an error by ID",
		Usage:   "<error ID>",

		CustomPermissions: admins,
		Command:           bot.error,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "changelog",
		Summary: "Show a list of terms added since `date`.\n`date` must be in `yyyy-mm-dd` format.",
		Usage:   "<channel> <date>",

		CustomPermissions: admins,
		Command:           bot.changelog,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "update-tags",
		Summary: "Bulk update a list of term's tags. Input in CSV format",

		CustomPermissions: admins,
		Command:           bot.updateTags,
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
		Command:           bot.importFromMessage,
	})

	a.AddSubcommand(&bcr.Command{
		Name:    "upload",
		Summary: "Upload a file",

		CustomPermissions: directors,
		Command:           bot.upload,
	})

	contributors := a.AddSubcommand(&bcr.Command{
		Name:              "contributor",
		Summary:           "Add contributors",
		CustomPermissions: directors,
		Command: func(ctx *bcr.Context) (err error) {
			return ctx.Help([]string{"admin", "contributor"})
		},
	})

	contributors.AddSubcommand(&bcr.Command{
		Name:              "category",
		Summary:           "List contributor categories",
		CustomPermissions: directors,
		Command:           bot.listContributorCategories,
	}).AddSubcommand(&bcr.Command{
		Name:              "add",
		Summary:           "Add a contributor category",
		Usage:             "<name> [role]",
		Args:              bcr.MinArgs(1),
		CustomPermissions: admins,
		Command:           bot.addContributorCategory,
	})

	contributors.AddSubcommand(&bcr.Command{
		Name:              "add",
		Summary:           "Add a contributor",
		Usage:             "<user> <category>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: directors,
		Command:           bot.addContributor,
	})

	contributors.AddSubcommand(&bcr.Command{
		Name:              "override",
		Summary:           "Override a contributor's name",
		Usage:             "<user> <new name|-clear>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: directors,
		Command:           bot.overrideContributor,
	})

	contributors.AddSubcommand(&bcr.Command{
		Name:      "import",
		Summary:   "Import all existing contributors (by role)",
		OwnerOnly: true,
		Command:   bot.allContributors,
	})

	i := bot.Router.AddCommand(bot.Router.AliasMust("ai", nil, []string{"admin", "import"}, nil))
	i.Args = bcr.MinArgs(1)

	bot.Router.AddHandler(bot.guildCreate)
	bot.Router.AddHandler(bot.guildDelete)

	bot.Router.ShardManager.ForEach(func(s shard.Shard) {
		state := s.(*state.State)

		state.AddHandler(func(_ *gateway.ReadyEvent) {
			var o sync.Once
			o.Do(func() {
				go bot.setStatusLoop(state)
			})
		})
	})

	auditlog.Init(b, directors)
	out = append(out, a)
	return "Bot admin commands", out
}
