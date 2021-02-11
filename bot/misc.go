package bot

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// UpdateStatus sets the bot's status to the given string
func (bot *Bot) UpdateStatus(name string, s gateway.Status) (err error) {
	return bot.Router.Session.Gateway.UpdateStatus(gateway.UpdateStatusData{
		Status: s,
		Activities: &[]discord.Activity{{
			Name: name,
		}},
	})
}
