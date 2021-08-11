package search

import (
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
)

func (c *commands) random(ctx bcr.Contexter) (err error) {
	catName := ctx.GetStringFlag("category")
	if catName == "" {
		if v, ok := ctx.(*bcr.Context); ok {
			catName = strings.Join(v.Args, " ")
		}
	}

	ignore := strings.Split(ctx.GetStringFlag("ignore"), ",")
	for i := range ignore {
		ignore[i] = strings.ToLower(strings.TrimSpace(ignore[i]))
	}

	// if theres arguments, try a category
	// returns true if it found a category
	if catName != "" {
		b, err := c.randomCategory(ctx, catName, ignore)
		if b || err != nil {
			return err
		}
	}

	// grab a random term
	t, err := c.DB.RandomTerm(ignore)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return ctx.SendEphemeral("No terms found! Are you sure you're not excluding every possible term?")
		}
		return c.DB.InternalError(ctx, err)
	}

	// send the random term
	_, err = ctx.Send("", c.DB.TermEmbed(t))
	return
}

func (c *commands) randomCategory(ctx bcr.Contexter, catName string, ignore []string) (b bool, err error) {
	cat, err := c.DB.CategoryID(catName)
	if err != nil {
		// dont bother to check if its a category not found error or not, just return nil
		return false, nil
	}

	t, err := c.DB.RandomTermCategory(cat, ignore)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			err = ctx.SendEphemeral("No terms found! Are you sure you're not excluding every possible term?")
			return true, err
		}
		return true, c.DB.InternalError(ctx, err)
	}

	err = ctx.SendX("", c.DB.TermEmbed(t))
	return true, err
}
