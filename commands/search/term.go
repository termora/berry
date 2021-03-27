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
		_, err = ctx.Send("You didn't give a term name or ID.", nil)
		return
	}

	var term *db.Term

	id, err := strconv.Atoi(ctx.RawArgs)
	if err == nil {
		term, err = c.DB.GetTerm(id)
		if err != nil {
			if errors.Cause(err) == pgx.ErrNoRows {
				_, err = ctx.Sendf("❌ No term with that ID found.")
				return
			}
			return c.DB.InternalError(ctx, err)
		}
	} else {
		term, err = c.DB.TermName(ctx.RawArgs)
		if err != nil && errors.Cause(err) != pgx.ErrNoRows {
			return c.DB.InternalError(ctx, err)
		} else if err == nil {
			goto found
		}

		{
			terms, err := c.DB.Search(ctx.RawArgs, 1)
			if err != nil {
				return c.DB.InternalError(ctx, err)
			}
			if len(terms) == 0 {
				_, err = ctx.Sendf("No term found.")
			}

			term = terms[0]
		}
	}

found:
	_, err = ctx.Send("", term.TermEmbed(c.Config.TermBaseURL()))
	return
}
