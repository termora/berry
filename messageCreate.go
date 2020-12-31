package main

import (
	"github.com/Starshine113/berry/structs"
	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type messageCreate struct {
	r     *crouter.Router
	c     *structs.BotConfig
	sugar *zap.SugaredLogger
}

func (mc *messageCreate) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// if message was sent by a bot return, unless it's in the list of allowed bots
	if m.Author.Bot && !inSlice(mc.c.Bot.AllowedBots, m.Author.ID) {
		return
	}

	// get context
	ctx, err := mc.r.Context(m)
	if err != nil {
		sugar.Error("Error creating context:", err)
		return
	}

	// check if the message might be a command
	if ctx.MatchPrefix() {
		mc.r.Execute(ctx)
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
