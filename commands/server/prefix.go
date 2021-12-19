package server

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) prefixes(ctx *bcr.Context) (err error) {
	prefixes := append(bot.Router.Prefixes, bot.PrefixesFor(ctx.Message.GuildID)...)

	// remove the first prefix, as the first two prefixes show up identical in the client
	prefixes = prefixes[1:]

	_, err = ctx.Send("", discord.Embed{
		Title:       "Prefixes",
		Description: strings.Join(prefixes, "\n"),
		Color:       ctx.Router.EmbedColor,
	})
	return
}
