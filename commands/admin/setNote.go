package admin

import (
	"strconv"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) setNote(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided: need ID and note (or \"clear\"", nil)
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

	note := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if note == "clear" {
		note = ""
	}

	if len(note) > 1000 {
		_, err = ctx.Sendf("❌ The CW you gave is too long (%v > 1000 characters).", len(note))
		return
	}

	err = c.db.SetNote(t.ID, note)
	if err != nil {
		return c.db.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Updated note for %v.", id)
	return
}
