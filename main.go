package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/commands/admin"
	"github.com/Starshine113/berry/commands/search"
	"github.com/Starshine113/berry/commands/server"
	"github.com/Starshine113/berry/commands/static"
	"github.com/Starshine113/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
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
	r.EmbedColor = 0xe00d7a

	// set blacklist function
	r.BlacklistFunc = d.CtxInBlacklist

	// add the message create handler
	mc := &messageCreate{r: r, c: c, sugar: sugar}
	s.AddHandler(mc.messageCreate)

	// set status
	s.AddHandler(func(d *gateway.ReadyEvent) {
		st := fmt.Sprintf("%vhelp", c.Bot.Prefixes[0])

		if c.Bot.Website != "" {
			var w string
			u, err := url.Parse(c.Bot.Website)
			if err != nil {
				w = c.Bot.Website
			} else {
				w = u.Host
			}

			st += " | " + w
		}

		if err := s.Gateway.UpdateStatus(gateway.UpdateStatusData{
			Status: gateway.OnlineStatus,
			Activities: &[]discord.Activity{{
				Name: st,
			}},
		}); err != nil {
			sugar.Error("Error setting status:", err)
		}
	})

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
