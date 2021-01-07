package admin

import (
	"strconv"
	"strings"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/misc"
)

func (c *commands) setCW(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided: need ID and bitmask", nil)
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

	cw := strings.Join(ctx.Args[1:], " ")

	if len(cw) > 1000 {
		_, err = ctx.Sendf("❌ The CW you gave is too long (%v > 1000 characters).", len(cw))
		return
	}

	err = c.db.SetCW(t.ID, cw)
	if err != nil {
		c.sugar.Errorf("Error setting CW for %v: %v", id, err)
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	_, err = ctx.Sendf("Updated CW for %v.", id)
	return
}
