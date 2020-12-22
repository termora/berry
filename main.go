package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/termbot/commands/admin"
	"github.com/Starshine113/termbot/commands/search"
	"github.com/Starshine113/termbot/commands/server"
	"github.com/Starshine113/termbot/commands/static"
	"github.com/Starshine113/termbot/db"
	"github.com/bwmarrin/discordgo"
)

var sugar *zap.SugaredLogger

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.RedirectStdLog(logger)
	sugar = logger.Sugar()

	c := getConfig(sugar)

	d, err := db.Init(c, sugar)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	sugar.Info("Connected to database.")

	dg, err := discordgo.New("Bot " + c.Auth.Token)
	if err != nil {
		sugar.Fatalf("Error creating Discord session: %v", err)
	}

	// create the router
	r := crouter.NewRouter(dg, c.Bot.BotOwners, c.Bot.Prefixes)
	// set blacklist function
	r.Blacklist(d.CtxInBlacklist)

	// add the message create handler
	dg.AddHandler(r.MessageCreate)

	// set status on connect
	dg.AddHandler(func(s *discordgo.Session, _ *discordgo.Ready) {
		if err := s.UpdateStatus(0, fmt.Sprintf("%vhelp", c.Bot.Prefixes[0])); err != nil {
			sugar.Errorf("Error setting status: %v", err)
		}
	})

	// add static commands
	static.Init(c, r)

	// add term commands
	search.Init(d, sugar, r)

	// add server commands
	server.Init(d, r)

	// add admin commands
	admin.Init(d, c, r)

	// add intents
	dg.Identify.Intents = discordgo.MakeIntent(crouter.RequiredIntents)

	// open a connection to Discord
	err = dg.Open()
	if err != nil {
		panic(err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		dg.Close()
		sugar.Infof("Disconnected from Discord.")
		d.Pool.Close()
		sugar.Infof("Closed database connection.")
	}()

	sugar.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
