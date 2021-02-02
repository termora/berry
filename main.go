package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/diamondburned/arikawa/v2/state"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/bot"
	"github.com/starshine-sys/berry/commands/admin"
	"github.com/starshine-sys/berry/commands/search"
	"github.com/starshine-sys/berry/commands/server"
	"github.com/starshine-sys/berry/commands/static"
	"github.com/starshine-sys/berry/db"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.RedirectStdLog(logger)
	sugar := logger.Sugar()

	c := getConfig(sugar)

	// command-line flags, mostly sharding
	pflag.BoolVarP(&c.Debug, "debug", "d", false, "Debug logging")
	pflag.IntVarP(&c.Shard, "shard", "s", 0, "Shard number")
	pflag.IntVarP(&c.NumShards, "shard-count", "c", 1, "Number of shards")
	pflag.Parse()
	c.Sharded = c.NumShards != 1

	// create a Sentry config
	if c.UseSentry {
		err = sentry.Init(sentry.ClientOptions{
			Dsn: c.Auth.SentryURL,
		})
		if err != nil {
			sugar.Fatalf("sentry.Init: %s", err)
		}
		sugar.Infof("Initialised Sentry")
		// defer this to flush buffered events
		defer sentry.Flush(2 * time.Second)
	}
	hub := sentry.CurrentHub()
	if !c.UseSentry {
		hub = nil
	}

	// connect to the database
	d, err := db.Init(c.Auth.DatabaseURL, sugar)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	d.SetSentry(hub)
	d.Config = c
	sugar.Info("Connected to database.")

	// create a new state
	s, err := state.NewWithIntents("Bot "+c.Auth.Token, bcr.RequiredIntents)
	if err != nil {
		log.Fatalln("Error creating state:", err)
	}

	// if the bot is sharded, set the number and count
	if c.Sharded {
		s.Gateway.Identifier.SetShard(c.Shard, c.NumShards)
	}

	// create a new router and set the default embed colour
	owners := make([]string, 0)
	for _, u := range c.Bot.BotOwners {
		owners = append(owners, u.String())
	}
	r := bcr.NewRouter(s, owners, c.Bot.Prefixes)
	r.EmbedColor = db.EmbedColour

	// set blacklist function
	r.BlacklistFunc = d.CtxInBlacklist

	// create the bot instance
	bot := bot.New(sugar, c, r, d, hub)
	// add search commands
	bot.Add(search.Init)
	// add static commands
	bot.Add(static.Init)
	// add server commands
	bot.Add(server.Init)
	// add admin commands
	bot.Add(admin.Init)

	// open a connection to Discord
	if err = s.Open(); err != nil {
		sugar.Fatal("Failed to connect:", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		s.Close()
		sugar.Infof("Disconnected from Discord.")
		d.Pool.Close()
		sugar.Infof("Closed database connection.")
	}()

	sugar.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	botUser, _ := s.Me()
	sugar.Infof("User: %v#%v (%v)", botUser.Username, botUser.Discriminator, botUser.ID)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
