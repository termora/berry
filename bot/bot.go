// Package bot contains the bot's core functionality.
package bot

import (
	"sort"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
	"github.com/starshine-sys/berry/structs"
	"go.uber.org/zap"
)

// Bot ...
type Bot struct {
	Sugar  *zap.SugaredLogger
	Config *structs.BotConfig
	Router *bcr.Router
	DB     *db.Db

	Modules []Module
}

// Module ...
type Module interface {
	String() string
	Commands() []*bcr.Command
}

// AddFunc ...
type AddFunc func(*Bot) (string, []*bcr.Command)

// New creates a new instance of Bot
func New(s *zap.SugaredLogger, config *structs.BotConfig, r *bcr.Router, db *db.Db) *Bot {
	b := &Bot{
		Sugar:  s,
		Config: config,
		Router: r,
		DB:     db,
	}

	b.Router.Session.AddHandler(b.MessageCreate)
	return b
}

// Add adds a module to the bot
func (bot *Bot) Add(f AddFunc) {
	m, c := f(bot)
	sort.Sort(bcr.Commands(c))
	bot.Modules = append(bot.Modules, &botModule{
		name:     m,
		commands: c,
	})
}
