package admin

import (
	"strconv"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/misc"
)

func (c *commands) delTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.Send("No term ID provided.", nil)
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

	m, err := ctx.Send("Are you sure you want to delete this term? React with ✅ to delete it, or with ❌ to cancel.", t.TermEmbed(""))
	if err != nil {
		return err
	}

	ctx.AddYesNoHandler(*m, ctx.Author.ID, func(ctx *bcr.Context) {
		err = c.db.RemoveTerm(id)
		if err != nil {
			c.sugar.Error("Error removing term:", err)
			ctx.Send(misc.InternalError, nil)
			return
		}
	}, func(ctx *bcr.Context) {
		ctx.Send("Cancelled.", nil)
		return
	})

	return nil
}
