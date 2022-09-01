package server

import (
	"strings"

	"github.com/starshine-sys/bcr"

	"github.com/termora/berry/db"
)

func (bot *Bot) blacklistAdd(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.")
		return err
	}

	// parse all channels passed to the command
	channels, n := ctx.GreedyChannelParser(ctx.Args)
	if n == 0 {
		_, err = ctx.Send("None of the channels you gave were found.")
		return
	}

	ch := make([]string, 0)
	for _, c := range channels {
		for _, cID := range ch {
			if cID == c.ID.String() {
				continue
			}
		}
		ch = append(ch, c.ID.String())
	}

	err = bot.DB.AddToBlacklist(ctx.Message.GuildID.String(), ch)
	if err != nil {
		if err == db.ErrorAlreadyBlacklisted {
			_, err = ctx.Send("One or more channels is already blacklisted.")
			return err
		}

		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Added %v to the blacklist.", strings.Join(mapString(ch, func(s string) string { return "<#" + s + ">" }), ", "))
	return
}

func (bot *Bot) blacklistRemove(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.")
		return err
	}

	ch, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		if err == bcr.ErrChannelNotFound {
			_, err = ctx.Send("The channel you gave was not found.")
			return err
		}

		return bot.DB.InternalError(ctx, err)
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = bot.DB.RemoveFromBlacklist(ctx.Message.GuildID.String(), ch.ID.String())
	if err != nil {
		if err == db.ErrorNotBlacklisted {
			_, err = ctx.Send("That channel isn't blacklisted.")
			return err
		}

		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Removed %v from the blacklist.", ch.Mention())
	return
}

// no generics :pensive:
func mapString(s []string, f func(string) string) []string {
	out := make([]string, 0)
	for _, e := range s {
		out = append(out, f(e))
	}
	return out
}
