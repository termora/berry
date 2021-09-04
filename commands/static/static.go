package static

import (
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
	"github.com/termora/berry/bot/cc"
	"github.com/termora/berry/commands/static/rpc"
	"google.golang.org/grpc"
)

// Commands ...
type Commands struct {
	*bot.Bot

	start time.Time

	memberMu             sync.RWMutex
	SupportServerMembers map[discord.UserID]discord.Member
	// is set to true once we receive the first GuildMembersChunk event
	guildMembersChunked bool

	rpc.UnimplementedGuildMemberServiceServer
}

// Init ...
func Init(bot *bot.Bot) (m string, o []*bcr.Command) {
	c := &Commands{
		Bot:                  bot,
		start:                time.Now().UTC(),
		SupportServerMembers: map[discord.UserID]discord.Member{},
	}

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "ping",

		Summary:  "Check the bot's message latency",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.ping,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name:    "stats",
		Aliases: []string{"about"},

		Summary:  "Show some statistics about the bot!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.about,
		Options:       &[]discord.CommandOption{},
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "feedback",

		Summary:  "Send feedback to the developers!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.feedback,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "credits",

		Summary:  "A list of people who helped create the bot",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.credits,
	}))

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name:    "hello",
		Aliases: []string{"Hi"},

		Summary:  "Say hi!",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       c.hello,
	}))

	help := bot.Router.AddCommand(&bcr.Command{
		Name:    "help",
		Aliases: []string{"h"},

		Summary:     "Show info about how to use the bot",
		Description: "Show info about how to use the bot. If a command name is given as an argument, show the help for that command.",
		Usage:       "[command]",
		Cooldown:    1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.help,
		Options:       &[]discord.CommandOption{},
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "permissions",
		Aliases: []string{"perms"},

		Summary:  "Show an explanation of the permissions the bot needs",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.perms,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "privacy",
		Aliases: []string{"privacy-policy"},

		Summary:  "Show an explanation of the data the bot collects (and what it doesn't)",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.privacy,
	})

	help.AddSubcommand(&bcr.Command{
		Name: "autopost",

		Summary:  "Show how to set up automatically posting terms in a channel",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.autopost,
	})

	help.AddSubcommand(&bcr.Command{
		Name:    "commands",
		Aliases: []string{"cmds"},

		Summary:  "Show a list of all commands",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		Command:       bot.CommandList,
	})

	o = append(o, bot.Router.AddCommand(&bcr.Command{
		Name: "invite",

		Summary:  "Get an invite link",
		Cooldown: 1 * time.Second,

		Blacklistable: true,
		SlashCommand:  c.cmdInvite,
		Options:       &[]discord.CommandOption{},
	}))

	export := bot.Router.AddCommand(&bcr.Command{
		Name:    "export",
		Summary: "Export all terms in a DM",
		Usage:   "[--compress|-x]",

		Command: c.export,
	})

	export.AddSubcommand(&bcr.Command{
		Name:    "csv",
		Summary: "Export terms as a CSV file",

		Command: c.exportCSV,
	})

	export.AddSubcommand(&bcr.Command{
		Name:    "xlsx",
		Summary: "Export terms as a XLSX file",

		Command: c.exportXLSX,
	})

	o = append(o, help, export)

	// thing
	if _, err := os.Stat("custom_commands.json"); err != nil {
		return "Bot info commands", o
	}

	bytes, err := os.ReadFile("custom_commands.json")
	if err != nil {
		bot.Sugar.Fatalf("Error reading custom commands file: %v", err)
	}

	cmds, err := cc.ParseBytes(bytes)
	if err != nil {
		bot.Sugar.Fatalf("Error parsing custom commands file: %v", err)
	}

	for _, c := range cmds {
		o = append(o, bot.Router.AddCommand(c))
	}

	c.newRPCServer()

	return "Bot info commands", o
}

func (c *Commands) newRPCServer() {
	port := strings.TrimPrefix(c.Config.RPCPort, ":")
	if port == "" {
		port = "58952"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		c.Sugar.Errorf("failed to listen: %v", err)
		return
	}

	rpcs := grpc.NewServer()
	rpc.RegisterGuildMemberServiceServer(rpcs, c)

	c.Sugar.Infof("RPC server listening at %v", lis.Addr())
	go func() {
		for {
			err := rpcs.Serve(lis)
			if err != nil {
				c.Sugar.Errorf("Failed to serve RPC: %v", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()
}
