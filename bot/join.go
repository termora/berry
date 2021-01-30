package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// GuildCreate logs the bot joining a server, and creates a database entry if one doesn't exist
func (bot *Bot) GuildCreate(g *gateway.GuildCreateEvent) {
	// create the server if it doesn't exist
	exists, err := bot.DB.CreateServerIfNotExists(g.ID.String())
	// if the server exists, don't log the join
	if exists {
		return
	}
	if err != nil {
		bot.Sugar.Errorf("Error creating database entry for server: %v", err)
		return
	}

	bot.Sugar.Infof("Joined guild %v (%v, owned by %v).", g.Name, g.ID, g.OwnerID)

	// if there's no channel to log joins/leaves to, return
	if bot.Config.Bot.JoinLogChannel == 0 {
		return
	}

	_, err = bot.Router.Session.SendMessageComplex(bot.Config.Bot.JoinLogChannel, api.SendMessageData{
		Content:         fmt.Sprintf("<a:ablobjoin:803960435983515648> Joined guild **%v** (%v, owned by <@%v>/%v)", g.Name, g.ID, g.OwnerID, g.OwnerID),
		AllowedMentions: &api.AllowedMentions{Parse: nil},
	})
	if err != nil {
		bot.Sugar.Errorf("Error sending log message: %v", err)
	}
	return
}

// GuildDelete logs the bot leaving a server and deletes the database entry
func (bot *Bot) GuildDelete(g *gateway.GuildDeleteEvent) {
	// if the guild's just unavailable, return, we didn't leave it
	if g.Unavailable {
		return
	}
	bot.Sugar.Infof("Left guild %v.", g.ID)

	// delete the server's database entry
	err := bot.DB.DeleteServer(g.ID.String())
	if err != nil {
		bot.Sugar.Errorf("Error deleting database entry for %v: %v", g.ID, err)
	}

	// if there's no channel to log joins/leaves to, return
	if bot.Config.Bot.JoinLogChannel == 0 {
		return
	}

	_, err = bot.Router.Session.SendMessageComplex(bot.Config.Bot.JoinLogChannel, api.SendMessageData{
		Content:         fmt.Sprintf("<a:ablobleave:803960446251171901> Left guild **%v**", g.ID),
		AllowedMentions: &api.AllowedMentions{Parse: nil},
	})
	if err != nil {
		bot.Sugar.Errorf("Error sending log message: %v", err)
	}
	return
}
