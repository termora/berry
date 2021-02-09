package search

import (
	"github.com/starshine-sys/bcr"
)

func (c *commands) random(ctx *bcr.Context) (err error) {
	// if theres arguments, try a category
	// returns true if it found a category
	if len(ctx.Args) > 0 {
		b, err := c.randomCategory(ctx)
		if b || err != nil {
			return err
		}
	}

	// grab a random term
	t, err := c.DB.RandomTerm()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	// send the random term
	_, err = ctx.Send("", t.TermEmbed(c.Config.Bot.TermBaseURL))
	return
}

func (c *commands) randomCategory(ctx *bcr.Context) (b bool, err error) {
	cat, err := c.DB.CategoryID(ctx.RawArgs)
	if err != nil {
		// dont bother to check if its a category not found error or not, just return nil
		return false, nil
	}

	t, err := c.DB.RandomTermCategory(cat)
	if err != nil {
		return true, c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("", t.TermEmbed(c.Config.Bot.TermBaseURL))
	return true, err
}
