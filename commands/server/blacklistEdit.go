package server

import (
	"strings"

	"github.com/Starshine113/berry/db"

	"github.com/Starshine113/crouter"
)

func (c *commands) blacklistAdd(ctx *crouter.Ctx) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	ch, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		return ctx.CommandError(err)
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = c.db.AddToBlacklist(ctx.Message.GuildID, ch.ID)
	if err != nil {
		if err == db.ErrorAlreadyBlacklisted {
			_, err = ctx.Send("That channel is already blacklisted.")
			return err
		}
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("Added %v to the blacklist.", ch.Mention())
	return
}

func (c *commands) blacklistRemove(ctx *crouter.Ctx) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	ch, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		return ctx.CommandError(err)
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = c.db.RemoveFromBlacklist(ctx.Message.GuildID, ch.ID)
	if err != nil {
		if err == db.ErrorNotBlacklisted {
			_, err = ctx.Send("That channel isn't blacklisted.")
			return err
		}
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("Removed %v from the blacklist.", ch.Mention())
	return
}
