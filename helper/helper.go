package helper

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
)

// Helper ...
type Helper struct {
	*session.Session

	GuildID discord.GuildID
	DB      *db.DB
}

const intents = gateway.IntentGuildMembers | gateway.IntentGuildMessages

// New creates a new Helper, adds the required intents and event handlers, and opens the connection.
func New(token string, id discord.GuildID, db *db.DB) (*Helper, error) {
	s := session.NewWithIntents("Bot "+token, intents)

	h := &Helper{
		Session: s,
		DB:      db,
		GuildID: id,
	}

	h.AddHandler(h.memberUpdate)

	err := h.Open(context.Background())
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (bot *Helper) memberUpdate(ev *gateway.GuildMemberUpdateEvent) {
	if ev.GuildID != bot.GuildID {
		return
	}

	name := ev.User.Username
	if ev.Nick != "" {
		name = ev.Nick
	}

	err := bot.DB.UpdateContributorName(ev.User.ID, name)
	if err != nil {
		log.Errorf("Error updating name for contributor: %v", err)
	}
}
