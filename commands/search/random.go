package search

import (
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
)

func (c *commands) random(ctx *bcr.Context) (err error) {
	ignore, _ := ctx.Flags.GetStringSlice("ignore-tags")
	for i := range ignore {
		ignore[i] = strings.ToLower(strings.TrimSpace(ignore[i]))
	}

	// if theres arguments, try a category
	// returns true if it found a category
	if len(ctx.Args) > 0 {
		b, err := c.randomCategory(ctx, ignore)
		if b || err != nil {
			return err
		}
	}

	// grab a random term
	t, err := c.DB.RandomTerm(ignore)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("No terms found! Are you sure you're not excluding every possible term?", nil)
			return
		}
		return c.DB.InternalError(ctx, err)
	}

	// send the random term
	_, err = ctx.Send("", t.TermEmbed(c.Config.TermBaseURL()))
	return
}

func (c *commands) randomCategory(ctx *bcr.Context, ignore []string) (b bool, err error) {
	cat, err := c.DB.CategoryID(ctx.RawArgs)
	if err != nil {
		// dont bother to check if its a category not found error or not, just return nil
		return false, nil
	}

	t, err := c.DB.RandomTermCategory(cat, ignore)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("No terms found! Are you sure you're not excluding every possible term?", nil)
			return true, nil
		}
		return true, c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("", t.TermEmbed(c.Config.TermBaseURL()))
	return true, err
}
