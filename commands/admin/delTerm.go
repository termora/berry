package admin

import (
	"strconv"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) delTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.Send("No term ID provided.", nil)
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	t, err := c.DB.GetTerm(id)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	m, err := ctx.Send("Are you sure you want to delete this term? React with ✅ to delete it, or with ❌ to cancel.", t.TermEmbed(""))
	if err != nil {
		return err
	}

	// confirm deleting the term
	if yes, timeout := ctx.YesNoHandler(*m, ctx.Author.ID); !yes || timeout {
		ctx.Send("Cancelled.", nil)
		return
	}

	err = c.DB.RemoveTerm(id)
	if err != nil {
		c.Sugar.Error("Error removing term:", err)
		c.DB.InternalError(ctx, err)
		return
	}
	_, err = ctx.Send("✅ Term deleted.", nil)
	if err != nil {
		c.Sugar.Error("Error sending message:", err)
	}

	return nil
}
