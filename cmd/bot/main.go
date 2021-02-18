package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	bcrbot "github.com/starshine-sys/bcr/bot"
	"github.com/starshine-sys/berry/bot"
	"github.com/starshine-sys/berry/commands/admin"
	"github.com/starshine-sys/berry/commands/pronouns"
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
	r, err := bcr.NewWithState(c.Auth.Token, c.Bot.BotOwners, c.Bot.Prefixes)
	if err != nil {
		sugar.Fatalf("Error creating router: %v", err)
	}

	// if the bot is sharded, set the number and count
	if c.Sharded {
		r.Session.Gateway.Identifier.SetShard(c.Shard, c.NumShards)
	}

	// set the default embed colour and blacklist function
	r.EmbedColor = db.EmbedColour
	r.BlacklistFunc = d.CtxInBlacklist

	// create the bot instance
	bot := bot.New(
		bcrbot.NewWithRouter(r),
		sugar, c, d, hub)
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

	// open a connection to Discord
	if err = r.Session.Open(); err != nil {
		sugar.Fatal("Failed to connect:", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		r.Session.Close()
		sugar.Infof("Disconnected from Discord.")
		d.Pool.Close()
		sugar.Infof("Closed database connection.")
	}()

	sugar.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	botUser, _ := r.Session.Me()
	sugar.Infof("User: %v#%v (%v)", botUser.Username, botUser.Discriminator, botUser.ID)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
