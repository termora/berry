package admin

import (
	"strconv"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) setCW(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided: need ID and CW (or \"clear\")")
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	t, err := bot.DB.GetTerm(id)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	// h
	// didn't we have a helper function for this??? oh well
	cw := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if cw == "clear" || cw == "-clear" {
		cw = ""
	}

	// if it's too long, return
	if len(cw) > 1000 {
		_, err = ctx.Sendf("❌ The CW you gave is too long (%v > 1000 characters).", len(cw))
		return
	}

	err = bot.DB.SetCW(t.ID, cw)
	if err != nil {
		bot.Log.Errorf("Error setting CW for %v: %v", id, err)
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Updated CW for %v.", id)
	return
}
