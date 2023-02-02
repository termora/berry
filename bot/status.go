package bot

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/termora/berry/common/log"
)

type statusThing struct {
	state    *state.State
	stop     chan struct{}
	sentStop bool
}

func (s *statusThing) reset() {
	if s == nil {
		return
	}

	s.sentStop = true
	s.stop <- struct{}{}
}

func (bot *Bot) ready(ev *gateway.ReadyEvent) {
	if len(bot.statuses) <= ev.Shard.ShardID() {
		log.Fatalf("shard ID is out of range (len: %d, id: %d)", len(bot.statuses), ev.Shard.ShardID())
	}

	s := bot.statuses[ev.Shard.ShardID()]
	if s == nil {
		state := bot.Router.ShardManager.Shard(ev.Shard.ShardID()).(*state.State)

		s = &statusThing{
			state: state,
			stop:  make(chan struct{}),
		}
		bot.statuses[ev.Shard.ShardID()] = s
	} else {
		if !s.sentStop {
			s.stop <- struct{}{}
		}
	}

	go bot.statusLoop(s)
}

func (bot *Bot) statusLoop(s *statusThing) {
	showGuildCount := bot.Config.Bot.ShowGuildCount
	tick := time.NewTicker(2 * time.Minute)

	// TODO: fix status
	log.Debugf("TEMP: not updating status, cancelling loop")

	return

	for {
		select {
		case <-s.stop:
			tick.Stop()
			return
		case <-tick.C:
			status := "/help"
			// add term count to status
			if bot.Config.Bot.ShowTermCount {
				status = fmt.Sprintf("%v | %v terms", status, bot.DB.TermCount())
			}

			// add website or guild count
			if showGuildCount {
				bot.GuildsMu.Lock()
				count := len(bot.Guilds)
				bot.GuildsMu.Unlock()

				status = fmt.Sprintf("%v | in %v servers", status, count)
			} else if bot.Config.Bot.Website != "" {
				status = fmt.Sprintf("%v | %v", status, urlParse(bot.Config.Bot.Website))
			}

			if err := s.state.Gateway().Send(context.Background(), &gateway.UpdatePresenceCommand{
				Status: discord.OnlineStatus,
				Activities: []discord.Activity{{
					Name: status,
					Type: discord.GameActivity,
				}},
			}); err != nil {
				log.Error("Error setting status:", err)
			}

			if bot.Config.Bot.ShowGuildCount {
				// switch status on the next loop
				showGuildCount = !showGuildCount
			}
		}
	}
}

func urlParse(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}
