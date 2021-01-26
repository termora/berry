package admin

import (
	"fmt"
	"net/url"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
)

// this is in admin to better integrate with the `guilds` admin command
func (c *Admin) setStatusLoop(s *state.State) {
	st := fmt.Sprintf("%vhelp", c.config.Bot.Prefixes[0])
	var guilds int
	countChan := make(chan int, 1)

	go c.guildCount(s, countChan)

	for {
		status := fmt.Sprintf("%v | %v", st, urlParse(c.config.Bot.Website))
		if c.config.Bot.Website == "" {
			status = st
		}

		if err := s.Gateway.UpdateStatus(gateway.UpdateStatusData{
			Status: gateway.OnlineStatus,
			Activities: &[]discord.Activity{{
				Name: status,
			}},
		}); err != nil {
			c.sugar.Error("Error setting status:", err)
		}

		// wait a minute to switch to other status
		time.Sleep(1 * time.Minute)

		select {
		case g := <-countChan:
			guilds = g
		default:
		}

		if err := s.Gateway.UpdateStatus(gateway.UpdateStatusData{
			Status: gateway.OnlineStatus,
			Activities: &[]discord.Activity{{
				Name: fmt.Sprintf("%v | in %v servers", st, guilds),
			}},
		}); err != nil {
			c.sugar.Error("Error setting status:", err)
		}

		// run once every minute
		time.Sleep(1 * time.Minute)
	}
}

func (c *Admin) guildCount(s *state.State, ch chan int) {
	for {
		// get number of guilds and send it over c
		g, err := s.Session.Client.Guilds(0)
		if err != nil {
			c.sugar.Error("Error getting guilds:", err)
			ch <- 0
		} else {
			ch <- len(g)
			c.guilds = g
		}

		// only run this once every hour
		time.Sleep(1 * time.Hour)
	}
}

func urlParse(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}
