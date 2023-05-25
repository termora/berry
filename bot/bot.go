// Package bot contains the bot's core functionality.
package bot

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/diamondburned/arikawa/v3/utils/httputil/httpdriver"
	"github.com/diamondburned/arikawa/v3/utils/ws"
	"github.com/getsentry/sentry-go"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/starshine-sys/bcr"
	bcrbot "github.com/starshine-sys/bcr/bot"
	"github.com/termora/berry/common"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
	"github.com/termora/berry/helper"
)

// Bot is the main bot struct
type Bot struct {
	*bcrbot.Bot
	Config common.Config
	DB     *db.DB

	Sentry    *sentry.Hub
	UseSentry bool

	Guilds   map[discord.GuildID]discord.Guild
	GuildsMu sync.Mutex

	Stats *StatsClient

	Helper       *helper.Helper
	StartStopLog *webhook.Client

	statuses []*statusThing
}

// New creates a new instance of Bot
func New(
	bot *bcrbot.Bot,
	config common.Config,
	db *db.DB, hub *sentry.Hub) *Bot {
	b := &Bot{
		Bot:       bot,
		Config:    config,
		DB:        db,
		Sentry:    hub,
		UseSentry: hub != nil,
		Guilds:    map[discord.GuildID]discord.Guild{},
	}

	if config.Bot.StartStopLog.ID.IsValid() && config.Bot.StartStopLog.Token != "" {
		b.StartStopLog = webhook.New(config.Bot.StartStopLog.ID, config.Bot.StartStopLog.Token)
	}

	// set the router's prefixer
	b.Router.Prefixer = b.Prefixer

	b.statuses = make([]*statusThing, b.Router.ShardManager.NumShards())
	id := 0
	// add the required handlers
	b.Router.ShardManager.ForEach(func(s shard.Shard) {
		sID := id
		state := s.(*state.State)

		state.PreHandler = handler.New()
		state.AddHandler(b.MessageCreate)
		state.AddHandler(b.InteractionCreate)
		state.AddHandler(b.GuildCreate)
		state.PreHandler.AddSyncHandler(b.GuildDelete)

		state.AddHandler(b.ready)

		state.AddHandler(func(ev *ws.CloseEvent) {
			if st := b.statuses[sID]; st != nil {
				st.reset()
			}

			log.Errorf("shard %v gateway closed with code %v: %v", sID, ev.Code, ev.Err)

			if b.StartStopLog != nil {
				err := b.StartStopLog.Execute(webhook.ExecuteData{
					Content: fmt.Sprintf("Shard %v gateway closed\n```Err: %v\nCode: %v\n```", sID, ev.Err, ev.Code),
				})
				if err != nil {
					log.Errorf("error sending log webhook: %v", err)
				}
			}

			if state.GatewayIsAlive() {
				log.Infof("attempting to close shard %v's connetion", sID)

				err := state.Close()
				if err != nil {
					log.Errorf("closing shard %v gateway: %v", sID, err)
				}
			}

			err := state.Open(context.Background())
			if err != nil {
				log.Errorf("reopening shard %v: %v", sID, err)
				if b.StartStopLog != nil {
					err := b.StartStopLog.Execute(webhook.ExecuteData{
						Content: fmt.Sprintf("Error reopening shard %v:\n```\n%v\n```", sID, err),
					})
					if err != nil {
						log.Errorf("error sending log webhook: %v", err)
					}
				}
			}
		})
		id++
	})

	// setup stats if metrics are enabled
	b.setupStats()

	if config.Bot.SupportToken != "" {
		h, err := helper.New(config.Bot.SupportToken, config.Bot.SupportGuildID, db)
		if err != nil {
			log.Errorf("Error creating helper: %v", err)
		}
		b.Helper = h
	} else {
		log.Warn("Helper token is NOT set! Admin commands may not work correctly.")
		log.Warn("If the main bot has the message content intent, set bot.support_token to the same token as the main bot.")
	}

	return b
}

// Add adds a module to the bot
func (bot *Bot) Add(f func(*Bot) (string, []*bcr.Command)) {
	m, c := f(bot)

	// sort the list of commands
	sort.Sort(bcr.Commands(c))

	// add the module
	bot.Modules = append(bot.Modules, &botModule{
		name:     m,
		commands: c,
	})
}

// Report reports an exception to Sentry if that's used, and the error is "our problem"
func (bot *Bot) Report(ctx *bcr.Context, err error) *sentry.EventID {
	if db.IsOurProblem(err) && bot.UseSentry {
		return bot.DB.CaptureError(ctx, err)
	}
	return nil
}

func (bot *Bot) guildCount() int {
	bot.GuildsMu.Lock()
	count := len(bot.Guilds)
	bot.GuildsMu.Unlock()
	return count
}

func (bot *Bot) setupStats() {
	if bot.Config.Bot.InfluxDB.URL != "" {
		log.Infof("Setting up InfluxDB client")

		bot.Stats = &StatsClient{
			Client:     influxdb2.NewClient(bot.Config.Bot.InfluxDB.URL, bot.Config.Bot.InfluxDB.Token).WriteAPI(bot.Config.Bot.InfluxDB.Org, bot.Config.Bot.InfluxDB.Bucket),
			guildCount: bot.guildCount,
		}

		bot.Router.ShardManager.ForEach(func(s shard.Shard) {
			state := s.(*state.State)

			state.Client.Client.OnResponse = append(state.Client.Client.OnResponse, func(httpdriver.Request, httpdriver.Response) error {
				go bot.Stats.IncAPICall()
				return nil
			})
		})

		bot.DB.IncFunc = bot.Stats.IncQuery

		go bot.Stats.submit()
	}
}
