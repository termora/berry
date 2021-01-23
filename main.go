package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/commands/admin"
	"github.com/starshine-sys/berry/commands/search"
	"github.com/starshine-sys/berry/commands/server"
	"github.com/starshine-sys/berry/commands/static"
	"github.com/starshine-sys/berry/db"
	"github.com/diamondburned/arikawa/v2/state"
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

	d, err := db.Init(c.Auth.DatabaseURL, sugar)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	sugar.Info("Connected to database.")

	s, err := state.NewWithIntents("Bot "+c.Auth.Token, bcr.RequiredIntents)
	if err != nil {
		log.Fatalln("Error creating state:", err)
	}

	r := bcr.NewRouter(s, c.Bot.BotOwners, c.Bot.Prefixes)
	r.EmbedColor = db.EmbedColour

	// set blacklist function
	r.BlacklistFunc = d.CtxInBlacklist

	// add the message create handler
	mc := &messageCreate{r: r, c: c, sugar: sugar}
	s.AddHandler(mc.messageCreate)

	// add static commands
	static.Init(c, d, sugar, r)

	// add term commands
	search.Init(d, c, sugar, r)

	// add server commands
	server.Init(d, r)

	// add admin commands
	admin.Init(d, sugar, c, r)

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
