package server

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) addPrefix(ctx *bcr.Context) (err error) {
	prefix := strings.ToLower(ctx.RawArgs)
	current := bot.PrefixesFor(ctx.Message.GuildID)

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a prefix to add.")
		return
	}

	if strings.Contains(prefix, bot.Router.Bot.ID.String()) {
		_, err = ctx.Send("Can't add mentioning the bot as a prefix.")
		return
	}

	for _, p := range current {
		if p == prefix {
			_, err = ctx.Sendf(":x: ``%v`` is already a prefix.", bcr.EscapeBackticks(ctx.RawArgs))
			return
		}
	}

	prefixes := append(current, prefix)

	con, cancel := bot.DB.Context()
	defer cancel()

	err = bot.DB.QueryRow(con, "update public.servers set prefixes = $1 where id = $2 returning prefixes", prefixes, ctx.Message.GuildID.String()).Scan(&prefixes)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("", discord.Embed{
		Title: "New prefixes",
		Description: strings.Join(
			append([]string{fmt.Sprintf("<@%v>", bot.Router.Bot.ID)}, prefixes...), "\n",
		),
		Color: bot.Router.EmbedColor,
	})
	return
}

func (bot *Bot) removePrefix(ctx *bcr.Context) (err error) {
	prefix := strings.ToLower(ctx.RawArgs)
	current := bot.PrefixesFor(ctx.Message.GuildID)

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a prefix to remove.")
		return
	}

	if strings.Contains(prefix, bot.Router.Bot.ID.String()) {
		_, err = ctx.Send("Can't remove mentioning the bot as a prefix.")
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

	con, cancel := bot.DB.Context()
	defer cancel()

	err = bot.DB.QueryRow(con, "update public.servers set prefixes = $1 where id = $2 returning prefixes", prefixes, ctx.Message.GuildID.String()).Scan(&prefixes)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("", discord.Embed{
		Title: "New prefixes",
		Description: strings.Join(
			append([]string{fmt.Sprintf("<@%v>", bot.Router.Bot.ID)}, prefixes...), "\n",
		),
		Color: bot.Router.EmbedColor,
	})
	return
}
