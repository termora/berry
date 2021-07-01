package server

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (c *commands) prefixes(ctx *bcr.Context) (err error) {
	prefixes := append(c.Router.Prefixes, c.PrefixesFor(ctx.Message.GuildID)...)

	// remove the first prefix, as the first two prefixes show up identical in the client
	prefixes = prefixes[1:]

	_, err = ctx.Send("", &discord.Embed{
		Title:       "Prefixes",
		Description: strings.Join(prefixes, "\n"),
		Color:       ctx.Router.EmbedColor,
	})
	return
}
