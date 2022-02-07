package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/termora/berry/common/log"
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
		log.Errorf("Error creating database entry for server: %v", err)
		return
	}

	log.Infof("Joined server %v (%v).", g.Name, g.ID)

	// if there's no channel to log joins/leaves to, return
	if bot.Config.Bot.JoinLogChannel == 0 {
		return
	}

	s, _ := bot.Router.StateFromGuildID(g.ID)

	_, err = s.SendMessageComplex(bot.Config.Bot.JoinLogChannel, api.SendMessageData{
		Content:         fmt.Sprintf("Joined new server **%v** (%v)", g.Name, g.ID),
		AllowedMentions: &api.AllowedMentions{Parse: nil},
	})
	if err != nil {
		log.Errorf("Error sending log message: %v", err)
	}
	return
}

// GuildDelete logs the bot leaving a server and deletes the database entry
func (bot *Bot) GuildDelete(g *gateway.GuildDeleteEvent) {
	// if the guild's just unavailable, return, we didn't leave it
	if g.Unavailable {
		return
	}

	// delete the server's database entry
	err := bot.DB.DeleteServer(g.ID.String())
	if err != nil {
		log.Errorf("Error deleting database entry for %v: %v", g.ID, err)
	}

	s, _ := bot.Router.StateFromGuildID(g.ID)

	guild, err := s.Guild(g.ID)
	if err != nil {
		// didn't find the guild, so just run this normally
		bot.guildDeleteNoState(s, g)
		return
	}

	// otherwise, use the cached guild
	log.Infof("Left server %v (%v)", guild.Name, guild.ID)

	// if there's no channel to log joins/leaves to, return
	if bot.Config.Bot.JoinLogChannel == 0 {
		return
	}

	_, err = s.SendMessageComplex(bot.Config.Bot.JoinLogChannel, api.SendMessageData{
		Content:         fmt.Sprintf("Left server **%v** (%v) :(", guild.Name, guild.ID),
		AllowedMentions: &api.AllowedMentions{Parse: nil},
	})
	if err != nil {
		log.Errorf("Error sending log message: %v", err)
	}
	return
}

// this is run if the left guild isn't found in the state
// which gives us almost no info, only the ID
func (bot *Bot) guildDeleteNoState(s *state.State, g *gateway.GuildDeleteEvent) {
	log.Infof("Left server %v.", g.ID)

	// if there's no channel to log joins/leaves to, return
	if bot.Config.Bot.JoinLogChannel == 0 {
		return
	}

	_, err := s.SendMessageComplex(bot.Config.Bot.JoinLogChannel, api.SendMessageData{
		Content:         fmt.Sprintf("Left server **%v** :(", g.ID),
		AllowedMentions: &api.AllowedMentions{Parse: nil},
	})
	if err != nil {
		log.Errorf("Error sending log message: %v", err)
	}
	return
}
