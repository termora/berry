package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	// start loop to update status every minute
	// :uhhh: arikawa doesn't have a way to add a handler that only runs once
	var o sync.Once
	s.AddHandler(func(d *gateway.ReadyEvent) {
		o.Do(func() {
			for {
				if err := s.Gateway.UpdateStatus(gateway.UpdateStatusData{
					Status: gateway.IdleStatus,
					Activities: &[]discord.Activity{{
						Name: fmt.Sprintf("%vhelp", c.Bot.Prefixes[0]),
					}},
				}); err != nil {
					sugar.Error("Error setting status:", err)
				}
				time.Sleep(time.Minute)
				// if a URL isn't set, just loop back immediately
				if c.Bot.Website == "" {
					continue
				}
				if err := s.Gateway.UpdateStatus(gateway.UpdateStatusData{
					Status: gateway.IdleStatus,
					Activities: &[]discord.Activity{{
						Name: fmt.Sprintf("%vhelp | %v", c.Bot.Prefixes[0], c.Bot.Website),
					}},
				}); err != nil {
					sugar.Error("Error setting status:", err)
				}
				time.Sleep(time.Minute)
			}
		})
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
