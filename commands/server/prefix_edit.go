package server

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (c *commands) addPrefix(ctx *bcr.Context) (err error) {
	prefix := strings.ToLower(ctx.RawArgs)
	current := c.PrefixesFor(ctx.Message.GuildID)

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a prefix to add.", nil)
		return
	}

	if strings.Contains(prefix, c.Router.Bot.ID.String()) {
		_, err = ctx.Send("Can't add mentioning the bot as a prefix.", nil)
		return
	}

	for _, p := range current {
		if p == prefix {
			_, err = ctx.Sendf(":x: ``%v`` is already a prefix.", bcr.EscapeBackticks(ctx.RawArgs))
			return
		}
	}

	prefixes := append(current, prefix)

	con, cancel := c.DB.Context()
	defer cancel()

	err = c.DB.Pool.QueryRow(con, "update public.servers set prefixes = $1 where id = $2 returning prefixes", prefixes, ctx.Message.GuildID.String()).Scan(&prefixes)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("", &discord.Embed{
		Title: "New prefixes",
		Description: strings.Join(
			append([]string{fmt.Sprintf("<@%v>", c.Router.Bot.ID)}, prefixes...), "\n",
		),
		Color: c.Router.EmbedColor,
	})
	return
}

func (c *commands) removePrefix(ctx *bcr.Context) (err error) {
	prefix := strings.ToLower(ctx.RawArgs)
	current := c.PrefixesFor(ctx.Message.GuildID)

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a prefix to remove.", nil)
		return
	}

	if strings.Contains(prefix, c.Router.Bot.ID.String()) {
		_, err = ctx.Send("Can't remove mentioning the bot as a prefix.", nil)
		return
	}

	var exists bool
	for _, p := range current {
		if p == prefix {
			exists = true
		}
	}
	if !exists {
		_, err = ctx.Sendf(":x: ``%v`` is not a prefix.", bcr.EscapeBackticks(ctx.RawArgs))
		return
	}

	// filter the prefixes
	prefixes := []string{}
	for _, p := range current {
		if p != prefix {
			prefixes = append(prefixes, p)
		}
	}

	con, cancel := c.DB.Context()
	defer cancel()

	err = c.DB.Pool.QueryRow(con, "update public.servers set prefixes = $1 where id = $2 returning prefixes", prefixes, ctx.Message.GuildID.String()).Scan(&prefixes)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("", &discord.Embed{
		Title: "New prefixes",
		Description: strings.Join(
			append([]string{fmt.Sprintf("<@%v>", c.Router.Bot.ID)}, prefixes...), "\n",
		),
		Color: c.Router.EmbedColor,
	})
	return
}
