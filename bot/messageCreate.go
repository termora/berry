package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/berry/db"
)

// MessageCreate is run when a message is created and handles commands
func (bot *Bot) MessageCreate(m *gateway.MessageCreateEvent) {
	var err error

	// debug logging, mostly for testing sharding *please don't use this in production i Beg*
	if bot.Config.Debug {
		bot.Sugar.Debugf("Message received from %v#%v (%v): %v", m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)
	}

	// defer panic handling
	defer func() {
		r := recover()
		if r != nil {
			bot.Sugar.Errorf("Caught panic in channel ID %v (user %v, guild %v): %v", m.ChannelID, m.Author.ID, m.GuildID, err)
		}
	}()

	// if the bot user isn't set yet, do that here
	// we can't do it when initialising the router because the connection to Discord will error
	if bot.Router.Bot == nil {
		err = bot.Router.SetBotUser()
		if err != nil {
			bot.Sugar.Error("Error setting bot user:", err)
			return
		}
		bot.Router.Prefixes = append(bot.Router.Prefixes, fmt.Sprintf("<@%v>", bot.Router.Bot.ID), fmt.Sprintf("<@!%v>", bot.Router.Bot.ID))
	}

	// if message was sent by a bot return, unless it's in the list of allowed bots
	if m.Author.Bot && !inSlice(bot.Config.Bot.AllowedBots, m.Author.ID) {
		return
	}

	// get context
	ctx, err := bot.Router.NewContext(m.Message)
	if err != nil {
		bot.Sugar.Error("Error creating context:", err)
		return
	}

	// check if the message might be a command
	if bot.Router.MatchPrefix(m.Content) {
		err = bot.Router.Execute(ctx)
		if err != nil {
			if db.IsOurProblem(err) {
				bot.Sentry.CaptureException(err)
			}
			bot.Sugar.Error(err)
		}
	}
}

func inSlice(slice []discord.UserID, s discord.UserID) bool {
	for _, i := range slice {
		if i == s {
			return true
		}
	}
	return false
}
