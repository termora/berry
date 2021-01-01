package admin

import (
	"strconv"
	"strings"

	"github.com/Starshine113/crouter"
)

func (c *commands) setCW(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		return ctx.CommandError(err)
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return ctx.CommandError(err)
	}

	t, err := c.db.GetTerm(id)
	if err != nil {
		return ctx.CommandError(err)
	}

	cw := strings.Join(ctx.Args[1:], " ")

	if len(cw) > 1000 {
		_, err = ctx.Sendf("❌ The CW you gave is too long (%v > 1000 characters).", len(cw))
		return
	}

	err = c.db.SetCW(t.ID, cw)
	if err != nil {
		c.sugar.Errorf("Error setting CW for %v: %v", id, err)
		return ctx.CommandError(err)
	}

	_, err = ctx.Sendf("Updated CW for %v.", id)
	return
}
