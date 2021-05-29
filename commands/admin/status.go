package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/state"
)

// this is in admin to better integrate with the `guilds` admin command
func (c *Admin) setStatusLoop(s *state.State) {
	st := fmt.Sprintf("%vhelp", c.Config.Bot.Prefixes[0])
	var guilds int
	countChan := make(chan int, 1)

	// spin off a function to fetch the guild count (well, actually fetch all guilds)
	// it's also used by `t!admin guilds`, which is why we run this even if the server count isn't shown in the bot's status
	go c.guildCount(s, countChan)

	for {
		// if something else set a static status, return
		select {
		case <-c.stopStatus:
			c.Sugar.Infof("Status loop stopped.")
			return
		default:
		}

		status := st
		// add term count to status
		if c.Config.Bot.ShowTermCount {
			status = fmt.Sprintf("%v | %v terms", st, c.DB.TermCount())
		}

		// add the website to the status, if it's not empty
		status = fmt.Sprintf("%v | %v", status, urlParse(c.Config.Bot.Website))
		if c.Config.Bot.Website == "" {
			status = st
		}
		// if the bot is sharded, also add the shard number to the status
		if c.Config.Sharded && c.Config.Bot.ShowShard {
			status = fmt.Sprintf("%v | shard %v/%v", status, s.Gateway.Identifier.Shard.ShardID(), s.Gateway.Identifier.Shard.NumShards())
		}

		if err := c.UpdateStatus(status, "online"); err != nil {
			c.Sugar.Error("Error setting status:", err)
		}

		// wait two minutes to switch to other status
		time.Sleep(2 * time.Minute)

		// if the guild count is disabled, loop immediately
		if !c.Config.Bot.ShowGuildCount {
			continue
		}

		// same as above--if a static status was set, return
		select {
		case <-c.stopStatus:
			c.Sugar.Infof("Status loop stopped.")
			return
		case g := <-countChan:
			guilds = g
		default:
		}

		status = st
		// add term count to status
		if c.Config.Bot.ShowTermCount {
			status = fmt.Sprintf("%v | %v terms", st, c.DB.TermCount())
		}

		status = fmt.Sprintf("%v | in %v servers", status, guilds)
		if c.Config.Sharded && c.Config.Bot.ShowShard {
			status = fmt.Sprintf("%v | shard %v", status, s.Gateway.Identifier.Shard.ShardID())
		}

		if err := c.UpdateStatus(status, "online"); err != nil {
			c.Sugar.Error("Error setting status:", err)
		}

		// run once every two minutes
		time.Sleep(2 * time.Minute)
	}
}

func (c *Admin) guildCount(s *state.State, ch chan int) {
	for {
		// get number of guilds and send it over c
		g, err := s.Session.Client.Guilds(0)
		if err != nil {
			c.Sugar.Error("Error getting guilds:", err)
			ch <- 0
		} else {
			ch <- len(g)
			// set the list of guilds in c, used for the `guilds` admin command
			c.guilds = g
			c.GuildCount = int64(len(g))

			// post guild count if needed
			if err = c.postGuildCount(len(g)); err != nil {
				c.Sugar.Errorf("Error posting guild count: %v", err)
			}
		}

		// only run this once every hour
		time.Sleep(1 * time.Hour)
	}
}

func (c *Admin) postGuildCount(count int) (err error) {
	u, err := c.Router.State.Me()
	if err != nil {
		return
	}

	if c.Config.BotLists.TopGG != "" {
		url := fmt.Sprintf("https://top.gg/api/bots/%v/stats", u.ID)

		c.Sugar.Infof("Posting guild count (%v) to top.gg's API", count)

		body := fmt.Sprintf(`{"server_count": %v}`, count)

		req, err := http.NewRequest("POST", url, strings.NewReader(body))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", c.Config.BotLists.TopGG)

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		c.Sugar.Infof("Posted guild count to top.gg's API")
	}

	if c.Config.BotLists.BotsGG != "" {
		url := fmt.Sprintf("https://discord.bots.gg/api/v1/bots/%v/stats", u.ID)

		c.Sugar.Infof("Posting guild count (%v) to discord.bots.gg's API", count)

		body := fmt.Sprintf(`{"guildCount": %v}`, count)

		req, err := http.NewRequest("POST", url, strings.NewReader(body))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", c.Config.BotLists.BotsGG)

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		c.Sugar.Infof("Posted guild count to discord.bots.gg's API")
	}

	return nil
}

func urlParse(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}
