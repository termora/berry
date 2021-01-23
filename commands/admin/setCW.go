package admin

import (
	"strconv"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (c *commands) setCW(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided: need ID and CW (or \"clear\"", nil)
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return c.db.InternalError(ctx, err)
	}

	t, err := c.db.GetTerm(id)
	if err != nil {
		return c.db.InternalError(ctx, err)
	}

	cw := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if cw == "clear" {
		cw = ""
	}

	if len(cw) > 1000 {
		_, err = ctx.Sendf("❌ The CW you gave is too long (%v > 1000 characters).", len(cw))
		return
	}

	err = c.db.SetCW(t.ID, cw)
	if err != nil {
		c.sugar.Errorf("Error setting CW for %v: %v", id, err)
		return c.db.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Updated CW for %v.", id)
	return
}
