package static

import (
	"sync"
	"time"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/structs"
	"go.uber.org/zap"
)

// Commands ...
type Commands struct {
	config *structs.BotConfig
	start  time.Time
	sugar  *zap.SugaredLogger
	db     *db.Db

	cmdCount int
	cmdMutex sync.RWMutex
}

// Init ...
func Init(conf *structs.BotConfig, d *db.Db, s *zap.SugaredLogger, r *bcr.Router) *Commands {
	c := &Commands{config: conf, start: time.Now(), sugar: s, db: d}
	r.AddCommand(&bcr.Command{
		Name: "ping",

		Summary:  "Check the bot's message latency",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.ping,
	})

	r.AddCommand(&bcr.Command{
		Name: "about",

		Summary:  "Some info about the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.about,
	})

	r.AddCommand(&bcr.Command{
		Name: "credits",

		Summary:  "A list of people who helped create the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.credits,
	})

	r.AddCommand(&bcr.Command{
		Name:    "hello",
		Aliases: []string{"Hi"},

		Summary:  "Say hi!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.hello,
	})

	r.AddCommand(&bcr.Command{
		Name: "help",

		Summary:  "Show info about how to use the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.help,
	})

	r.AddCommand(&bcr.Command{
		Name: "invite",

		Summary:  "Get an invite link",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.cmdInvite,
	})

	return c
}

// PostFunc logs when a command is used
func (c *Commands) PostFunc(ctx *bcr.Context) {
	c.sugar.Debugf("Command executed: `%v` (arguments %v) by %v (channel %v, guild %v)", ctx.Cmd.Name, ctx.Args, ctx.Author.ID, ctx.Channel.ID, ctx.Message.GuildID)
	c.cmdMutex.Lock()
	c.cmdCount++
	c.cmdMutex.Unlock()
}
