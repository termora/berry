package static

import (
	"sync"
	"time"

	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/structs"
	"github.com/Starshine113/crouter"
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
func Init(conf *structs.BotConfig, d *db.Db, s *zap.SugaredLogger, r *crouter.Router) *Commands {
	c := &Commands{config: conf, start: time.Now(), sugar: s, db: d}
	r.AddCommand(&crouter.Command{
		Name: "ping",

		Summary:  "Check the bot's message latency",
		Cooldown: 3 * time.Second,

		Blacklistable: true,
		Command:       c.ping,
	})

	r.AddCommand(&crouter.Command{
		Name: "about",

		Summary:  "Some info about the bot",
		Cooldown: 5 * time.Second,

		Blacklistable: true,
		Command:       c.about,
	})

	r.AddCommand(&crouter.Command{
		Name:    "hello",
		Aliases: []string{"Hi"},

		Summary:  "Say hi!",
		Cooldown: 3 * time.Second,

		Blacklistable: true,
		Command:       c.hello,
	})

	r.AddCommand(&crouter.Command{
		Name: "help",

		Summary:  "Show info about how to use the bot",
		Cooldown: 5 * time.Second,

		Blacklistable: true,
		Command:       c.help,
	})

	r.AddCommand(&crouter.Command{
		Name: "invite",

		Summary:  "Get an invite link",
		Cooldown: 5 * time.Second,

		Blacklistable: true,
		Command:       c.cmdInvite,
	})

	return c
}

// PostFunc logs when a command is used
func (c *Commands) PostFunc(ctx *crouter.Ctx) {
	c.sugar.Debugf("Command executed: `%v` (arguments %v) by %v (channel %v, guild %v)", ctx.Cmd.Name, ctx.Args, ctx.Author.ID, ctx.Channel.ID, ctx.Message.GuildID)
	c.cmdMutex.Lock()
	c.cmdCount++
	c.cmdMutex.Unlock()
}
