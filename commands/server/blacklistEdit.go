package server

import (
	"strings"

	"github.com/Starshine113/bcr"

	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/misc"
)

func (c *commands) blacklistAdd(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.", nil)
		return err
	}

	ch, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		if err == bcr.ErrChannelNotFound {
			_, err = ctx.Send("The channel you gave was not found.", nil)
			return err
		}

		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = c.db.AddToBlacklist(ctx.Message.GuildID.String(), ch.ID.String())
	if err != nil {
		if err == db.ErrorAlreadyBlacklisted {
			_, err = ctx.Send("That channel is already blacklisted.", nil)
			return err
		}

		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	_, err = ctx.Sendf("Added %v to the blacklist.", ch.Mention())
	return
}

func (c *commands) blacklistRemove(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.", nil)
		return err
	}

	ch, err := ctx.ParseChannel(strings.Join(ctx.Args, " "))
	if err != nil {
		if err == bcr.ErrChannelNotFound {
			_, err = ctx.Send("The channel you gave was not found.", nil)
			return err
		}

		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = c.db.RemoveFromBlacklist(ctx.Message.GuildID.String(), ch.ID.String())
	if err != nil {
		if err == db.ErrorNotBlacklisted {
			_, err = ctx.Send("That channel isn't blacklisted.", nil)
			return err
		}

		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	_, err = ctx.Sendf("Removed %v from the blacklist.", ch.Mention())
	return
}
