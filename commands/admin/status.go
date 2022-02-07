package admin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/termora/berry/common/log"
)

// this is in admin to better integrate with the `guilds` admin command
func (bot *Bot) setStatusLoop(s *state.State) {
	st := fmt.Sprintf("%vhelp", bot.Config.Bot.Prefixes[0])

	// spin off a function to fetch the guild count (well, actually fetch all guilds)
	// it's also used by `t!admin guilds`, which is why we run this even if the server count isn't shown in the bot's status
	if s.Ready().Shard.ShardID() == 0 {
		log.Debugf("Spawning guildCount function on shard 0")
		go bot.guildCount(s)
	}

	for {
		// if something else set a static status, return
		select {
		case <-bot.stopStatus:
			log.Infof("Status loop stopped.")
			return
		default:
		}

		status := st
		// add term count to status
		if bot.Config.Bot.ShowTermCount {
			status = fmt.Sprintf("%v | %v terms", st, bot.DB.TermCount())
		}

		// add the website to the status, if it's not empty
		if bot.Config.Bot.Website != "" {
			status = fmt.Sprintf("%v | %v", status, urlParse(bot.Config.Bot.Website))
		}
		// if the bot is sharded, also add the shard number to the status
		if bot.Router.ShardManager.NumShards() > 1 && bot.Config.Bot.ShowShard {
			status = fmt.Sprintf("%v | shard %v/%v", status, s.Ready().Shard.ShardID()+1, bot.Router.ShardManager.NumShards())
		}

		if err := s.Gateway().Send(context.Background(), &gateway.UpdatePresenceCommand{
			Status: discord.OnlineStatus,
			Activities: []discord.Activity{{
				Name: status,
				Type: discord.GameActivity,
			}},
		}); err != nil {
			log.Error("Error setting status:", err)
		}

		// wait two minutes to switch to other status
		time.Sleep(2 * time.Minute)

		// if the guild count is disabled, loop immediately
		if !bot.Config.Bot.ShowGuildCount {
			continue
		}

		// same as above--if a static status was set, return
		select {
		case <-bot.stopStatus:
			log.Infof("Status loop stopped.")
			return
		default:
		}

		status = st
		// add term count to status
		if bot.Config.Bot.ShowTermCount {
			status = fmt.Sprintf("%v | %v terms", st, bot.DB.TermCount())
		}

		bot.GuildsMu.Lock()
		count := len(bot.Guilds)
		bot.GuildsMu.Unlock()

		status = fmt.Sprintf("%v | in %v servers", status, count)

		// if the bot is sharded, also add the shard number to the status
		if bot.Router.ShardManager.NumShards() > 1 && bot.Config.Bot.ShowShard {
			status = fmt.Sprintf("%v | shard %v/%v", status, s.Ready().Shard.ShardID()+1, bot.Router.ShardManager.NumShards())
		}

		if err := s.Gateway().Send(context.Background(), &gateway.UpdatePresenceCommand{
			Status: discord.OnlineStatus,
			Activities: []discord.Activity{{
				Name: status,
				Type: discord.GameActivity,
			}},
		}); err != nil {
			log.Error("Error setting status:", err)
		}

		// run once every two minutes
		time.Sleep(2 * time.Minute)
	}
}

func (bot *Bot) guildCount(s *state.State) {
	time.Sleep(5 * time.Minute)

	for {
		// post guild count if needed
		bot.GuildsMu.Lock()
		count := len(bot.Guilds)
		bot.GuildsMu.Unlock()

		if err := bot.postGuildCount(s, count); err != nil {
			log.Errorf("Error posting guild count: %v", err)
		}

		// only run this once every hour
		time.Sleep(1 * time.Hour)
	}
}

func (bot *Bot) postGuildCount(s *state.State, count int) (err error) {
	u, err := s.Me()
	if err != nil {
		return
	}

	if bot.Config.BotLists.TopGG != "" {
		url := fmt.Sprintf("https://top.gg/api/bots/%v/stats", u.ID)

		log.Infof("Posting guild count (%v) to top.gg's API", count)

		body := fmt.Sprintf(`{"server_count": %v}`, count)

		req, err := http.NewRequest("POST", url, strings.NewReader(body))
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", bot.Config.BotLists.TopGG)

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		log.Infof("Posted guild count to top.gg's API")
	}

	if bot.Config.BotLists.BotsGG != "" {
		url := fmt.Sprintf("https://discord.bots.gg/api/v1/bots/%v/stats", u.ID)

		log.Infof("Posting guild count (%v) to discord.bots.gg's API", count)

		body := fmt.Sprintf(`{"guildCount": %v}`, count)

		req, err := http.NewRequest("POST", url, strings.NewReader(body))
		if err != nil {
			return err
		}

		site := bot.Config.Bot.Website
		if site == "" {
			site = bot.Config.Bot.Git
		}
		req.Header["User-Agent"] = []string{fmt.Sprintf("%v-%v/v6 (Arikawa; +%v) DBots/%v", u.Username, u.Discriminator, site, u.ID)}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", bot.Config.BotLists.BotsGG)

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		log.Infof("Posted guild count to discord.bots.gg's API")
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
