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

	help := r.AddCommand(&bcr.Command{
		Name: "help",

		Summary:     "Show info about how to use the bot",
		Description: "Show info about how to use the bot. If a command name is given as an argument, show the help for that command.",
		Usage:       "[command]",
		Cooldown:    1 * time.Second,

		Blacklistable: true,
		Command:       c.help,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "permissions",
		Aliases: []string{"perms"},

		Summary:  "Show an explanation of the permissions the bot needs",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.perms,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "privacy",
		Aliases: []string{"privacy-policy"},

		Summary:  "Show an explanation of the data the bot collects (and what it doesn't)",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.privacy,
	})

	help.AddSubcommand(&bcr.Command{
		Name: "autopost",

		Summary:  "Show how to set up automatically posting terms in a channel",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.autopost,
	})

	help.AddSubcommand(&bcr.Command{
		Name: "commands",

		Summary:  "Show a list of all commands",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.commandList,
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
