// Package bot contains the bot's core functionality.
package bot

import (
	"sort"

	"github.com/diamondburned/arikawa/v2/utils/handler"
	"github.com/getsentry/sentry-go"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
	"github.com/starshine-sys/berry/structs"
	"go.uber.org/zap"
)

// Bot is the main bot struct
type Bot struct {
	Sugar  *zap.SugaredLogger
	Config *structs.BotConfig
	Router *bcr.Router
	DB     *db.Db

	Modules []Module

	Sentry    *sentry.Hub
	UseSentry bool
}

// Module is a single module/category of commands
type Module interface {
	String() string
	Commands() []*bcr.Command
}

// New creates a new instance of Bot
func New(s *zap.SugaredLogger, config *structs.BotConfig, r *bcr.Router, db *db.Db, hub *sentry.Hub) *Bot {
	b := &Bot{
		Sugar:     s,
		Config:    config,
		Router:    r,
		DB:        db,
		Sentry:    hub,
		UseSentry: hub != nil,
	}

	// create a pre-handler
	b.Router.Session.PreHandler = handler.New()
	b.Router.Session.PreHandler.Synchronous = true

	// add the required handlers
	b.Router.Session.AddHandler(b.MessageCreate)
	b.Router.Session.AddHandler(b.GuildCreate)
	b.Router.Session.PreHandler.AddHandler(b.GuildDelete)
	return b
}

// Add adds a module to the bot
func (bot *Bot) Add(f func(*Bot) (string, []*bcr.Command)) {
	m, c := f(bot)

	// sort the list of commands
	sort.Sort(bcr.Commands(c))

	// add the module
	bot.Modules = append(bot.Modules, &botModule{
		name:     m,
		commands: c,
	})
}
