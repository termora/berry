package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Starshine113/berry/commands/admin"
	"github.com/Starshine113/berry/commands/search"
	"github.com/Starshine113/berry/commands/server"
	"github.com/Starshine113/berry/commands/static"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/crouter"
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

	d, err := db.Init(c.Auth.DatabaseURL, sugar)
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
	// set post-command log function
	r.PostFunc = postFunc

	// add the message create handler
	mc := &messageCreate{r: r, c: c, sugar: sugar}
	dg.AddHandler(mc.messageCreate)

	// start loop to update status every minute
	dg.AddHandlerOnce(func(s *discordgo.Session, _ *discordgo.Ready) {
		for {
			if err := s.UpdateStatus(0, fmt.Sprintf("%vhelp | in %v servers", c.Bot.Prefixes[0], len(s.State.Guilds))); err != nil {
				sugar.Errorf("Error setting status: %v", err)
			}
			time.Sleep(time.Minute)
			// if a URL isn't set, just loop back immediately
			if c.Bot.Website == "" {
				continue
			}
			if err := s.UpdateStatus(0, fmt.Sprintf("%vhelp | %v", c.Bot.Prefixes[0], c.Bot.Website)); err != nil {
				sugar.Errorf("Error setting status: %v", err)
			}
			time.Sleep(time.Minute)
		}
	})

	// add static commands
	static.Init(c, r)

	// add term commands
	search.Init(d, sugar, r)

	// add server commands
	server.Init(d, r)

	// add admin commands
	admin.Init(d, sugar, c, r)

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
