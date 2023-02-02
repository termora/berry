package bot

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/diamondburned/arikawa/v3/utils/ws"
	"github.com/getsentry/sentry-go"
	"github.com/urfave/cli/v2"

	"github.com/starshine-sys/bcr"
	bcrbot "github.com/starshine-sys/bcr/bot"
	"github.com/termora/berry/bot"
	"github.com/termora/berry/commands/admin"
	"github.com/termora/berry/commands/pronouns"
	"github.com/termora/berry/commands/search"
	"github.com/termora/berry/commands/server"
	"github.com/termora/berry/commands/static"
	"github.com/termora/berry/common"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
	dbsearch "github.com/termora/berry/db/search"
	"github.com/termora/berry/db/search/typesense"
)

var Command = &cli.Command{
	Name:   "bot",
	Usage:  "Run the bot",
	Action: run,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Debug logging",
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "noloop",
			Aliases: []string{"N"},
			Value:   false,
			Usage:   "Disable event loop that will kill bot after 5 minutes of no events",
		},
		&cli.BoolFlag{
			Name:  "more-debug",
			Value: false,
			Usage: "Even MORE debug logs (very spammy)",
		},
	},
}

func run(ctx *cli.Context) error {
	rand.Seed(time.Now().UnixNano())

	c := common.ReadConfig()

	if ctx.Bool("debug") {
		ws.WSDebug = log.Debug
		db.Debug = log.Debugf
	}

	// create a Sentry config
	if c.Core.UseSentry {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: c.Core.SentryURL,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
		log.Infof("Initialised Sentry")
		// defer this to flush buffered events
		defer sentry.Flush(2 * time.Second)
	}
	hub := sentry.CurrentHub()
	if !c.Core.UseSentry {
		hub = nil
	}

	// connect to the database
	d, err := db.Init(c.Core.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	d.SetSentry(hub)
	d.Config = c
	d.TermBaseURL = c.Bot.TermBaseURL()
	defer func() {
		d.Pool.Close()
		log.Infof("Closed database connection.")
	}()

	if c.Core.TypesenseURL != "" && c.Core.TypesenseKey != "" {
		d.Searcher, err = typesense.New(c.Core.TypesenseURL, c.Core.TypesenseKey, d.Pool)
		if err != nil {
			log.Fatalf("Error connecting to Typesense: %v", err)
		}
	}

	// sync terms
	terms, err := d.GetTerms(dbsearch.FlagSearchHidden)
	if err != nil {
		log.Fatalf("Couldn't fetch all terms: %v", err)
	}

	err = d.SyncTerms(terms)
	if err != nil {
		log.Fatalf("Couldn't synchronize terms: %v", err)
	}
	log.Info("Synchronized terms with search instance!")

	log.Info("Connected to database.")

	// create a new shard manager
	mgr, err := shard.NewIdentifiedManager(gateway.IdentifyCommand{
		Token:      "Bot " + c.Bot.Token,
		Properties: gateway.DefaultIdentity,
		Shard:      &gateway.Shard{0, 3}, // just hardcode the number of shards it's fiiiiiine
		Presence: &gateway.UpdatePresenceCommand{
			Status: discord.OnlineStatus,
			Activities: []discord.Activity{{
				Name: fmt.Sprintf("/help | %v", urlParse(c.Bot.Website)),
				Type: discord.GameActivity,
			}},
		},
	}, state.NewShardFunc(func(_ *shard.Manager, s *state.State) {
		s.AddIntents(bcr.RequiredIntents)
	}))
	if err != nil {
		log.Fatalf("Error creating shard manager: %v", err)
	}

	b := bcrbot.NewWithRouter(bcr.New(mgr, []string{}, c.Bot.Prefixes))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	b.Router.ShardManager.ForEach(func(s shard.Shard) {
		state := s.(*state.State)

		state.Cabinet.MessageStore = store.Noop

		state.AddHandler(func(err error) {
			log.Errorf("Gateway error: %v", err)
		})
	})

	b.Owner(c.Bot.BotOwners...)

	// set the default embed colour and blacklist function
	b.Router.EmbedColor = db.EmbedColour
	b.Router.BlacklistFunc = d.CtxInBlacklist

	// create the bot instance
	bot := bot.New(
		b, c, d, hub)
	// add search commands
	bot.Add(search.Init)
	// add pronoun commands
	bot.Add(pronouns.Init)
	// add static commands
	bot.Add(static.Init)
	// add server commands
	bot.Add(server.Init)
	// add admin commands
	bot.Add(admin.Init)

	state, _ := bot.Router.StateFromGuildID(0)
	botUser, _ := state.Me()
	bot.Router.Bot = botUser
	bot.Router.Prefixes = append(bot.Router.Prefixes, fmt.Sprintf("<@%v>", bot.Router.Bot.ID), fmt.Sprintf("<@!%v>", bot.Router.Bot.ID))

	// open a connection to Discord
	if err = bot.Start(context.Background()); err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		bot.Router.ShardManager.Close()
		log.Infof("Disconnected from Discord.")
	}()

	log.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")
	log.Infof("User: %v (%v)", botUser.Tag(), botUser.ID)

	if c.Bot.SlashEnabled {
		if len(c.Bot.SlashGuilds) > 0 {
			log.Infof("Syncing commands in %v...", c.Bot.SlashGuilds)
		} else {
			log.Info("Syncing slash commands...")
		}
		err = bot.Router.SyncCommands(c.Bot.SlashGuilds...)
		if err != nil {
			log.Errorf("Couldn't sync commands: %v", err)
		} else {
			log.Info("Synced commands!")
		}
	} else {
		log.Info("Not syncing slash commands.")
	}

	go timer()

	cctx, stop := signal.NotifyContext(ctx.Context, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	exitCh := make(chan struct{})
	if !ctx.Bool("noloop") {
		eventCh := make(chan interface{}, 100)

		go eventThing(ctx, eventCh, exitCh)

		bot.Router.AddHandler(eventCh)
	}

	select {
	case <-cctx.Done():
	case <-exitCh:
	}

	log.Infof("Interrupt signal received. Shutting down...")

	return nil
}

func timer() {
	t := time.Now().UTC()
	ch := time.Tick(10 * time.Minute)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	for {
		select {
		case <-ch:
			log.Debugf("Tick received, %s since last tick.", time.Since(t))
			t = time.Now().UTC()
		case <-ctx.Done():
			return
		}
	}
}

func eventThing(ctx *cli.Context, ch <-chan interface{}, out chan<- struct{}) {
	cctx, stop := signal.NotifyContext(ctx.Context, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	t := time.AfterFunc(5*time.Minute, func() {
		out <- struct{}{}
	})

	for {
		select {
		case ev := <-ch:
			if ctx.Bool("more-debug") {
				log.Debugf("Received event %s", reflect.ValueOf(ev).Elem().Type().Name())
			}
			t.Stop()
			t = time.AfterFunc(5*time.Minute, func() {
				out <- struct{}{}
			})
		case <-cctx.Done():
			// break if we're shutting down
			break
		}
	}
}

func urlParse(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}
