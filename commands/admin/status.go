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

	// spin off a function to fetch the guild count (well, actually fetch all guilds)
	// it's also used by `t!admin guilds`, which is why we run this even if the server count isn't shown in the bot's status
	go c.guildCount(s, countChan)

	for {
		// if something else set a static status, return
		select {
		case <-c.stopStatus:
			c.sugar.Infof("Status loop stopped.")
			return
		default:
		}

		// add the website to the status, if it's not empty
		status := fmt.Sprintf("%v | %v", st, urlParse(c.config.Bot.Website))
		if c.config.Bot.Website == "" {
			status = st
		}
		// if the bot is sharded, also add the shard number to the status
		if c.config.Sharded {
			status = fmt.Sprintf("%v | shard %v/%v", status, s.Gateway.Identifier.Shard.ShardID(), s.Gateway.Identifier.Shard.NumShards())
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

		// if the guild count is disabled, loop immediately
		if !c.config.Bot.ShowGuildCount {
			continue
		}

		// same as above--if a static status was set, return
		select {
		case <-c.stopStatus:
			c.sugar.Infof("Status loop stopped.")
			return
		case g := <-countChan:
			guilds = g
		default:
		}

		status = fmt.Sprintf("%v | in %v servers", st, guilds)
		if c.config.Sharded {
			status = fmt.Sprintf("%v | shard %v/%v", status, s.Gateway.Identifier.Shard.ShardID(), s.Gateway.Identifier.Shard.NumShards())
		}

		if err := s.Gateway.UpdateStatus(gateway.UpdateStatusData{
			Status: gateway.OnlineStatus,
			Activities: &[]discord.Activity{{
				Name: status,
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
			// set the list of guilds in c, used for the `guilds` admin command
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
