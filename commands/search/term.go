package search

import (
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/termora/berry/db"

	"github.com/starshine-sys/bcr"
)

func (c *commands) term(ctx *bcr.Context) (err error) {
	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a term name or ID.")
		return
	}

	var (
		exact bool
		term  *db.Term
	)

	id, err := strconv.Atoi(ctx.RawArgs)
	if err == nil {
		term, err = c.DB.GetTerm(id)
		if err != nil {
			if errors.Cause(err) == pgx.ErrNoRows {
				_, err = ctx.Sendf("No term with that ID found.")
				return
			}
			return c.DB.InternalError(ctx, err)
		}
		exact = true
	} else {
		terms, err := c.DB.TermName(ctx.RawArgs)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		} else if err == nil && len(terms) > 0 {
			if len(terms) > 1 {
				return c.search(ctx)
			}

			exact = true
			term = terms[0]
			goto found
		}

		{
			terms, err := c.DB.Search(ctx.RawArgs, 1, nil)
			if err != nil {
				return c.DB.InternalError(ctx, err)
			}
			if len(terms) == 0 {
				_, err = ctx.Sendf("No term found.")
				return err
			}

			term = terms[0]
		}
	}

found:
	m := ctx.NewMessage()

	if !exact {
		m = m.Content("I couldn't find a term exactly matching that name, but here's the closest match:")
	}

	e := term.TermEmbed(c.Config.TermBaseURL())

	_, err = m.Embeds(e).Send()
	return
}
