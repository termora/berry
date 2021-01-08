package admin

import (
	"strconv"
	"strings"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/misc"
)

func (c *commands) setNote(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided: need ID and note (or \"clear\"", nil)
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	t, err := c.db.GetTerm(id)
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
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
		c.sugar.Errorf("Error setting note for %v: %v", id, err)
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	_, err = ctx.Sendf("Updated note for %v.", id)
	return
}
