package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/gateway"
)

// MessageCreate ...
func (bot *Bot) MessageCreate(m *gateway.MessageCreateEvent) {
	var err error

	// defer panic handling
	defer func() {
		r := recover()
		if r != nil {
			bot.Sugar.Errorf("Caught panic in channel ID %v (user %v, guild %v): %v", m.ChannelID, m.Author.ID, m.GuildID, err)
		}
	}()

	if bot.Router.Bot == nil {
		err = bot.Router.SetBotUser()
		if err != nil {
			bot.Sugar.Error("Error setting bot user:", err)
			return
		}
		bot.Router.Prefixes = append(bot.Router.Prefixes, fmt.Sprintf("<@%v>", bot.Router.Bot.ID), fmt.Sprintf("<@!%v>", bot.Router.Bot.ID))
	}

	// if message was sent by a bot return, unless it's in the list of allowed bots
	if m.Author.Bot && !inSlice(bot.Config.Bot.AllowedBots, m.Author.ID.String()) {
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
		bot.Router.Execute(ctx)
	}
}

func inSlice(slice []string, s string) bool {
	for _, i := range slice {
		if i == s {
			return true
		}
	}
	return false
}
