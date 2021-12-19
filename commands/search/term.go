package search

import (
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/termora/berry/db"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) term(ctx *bcr.Context) (err error) {
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
		term, err = bot.DB.GetTerm(id)
		if err != nil {
			if errors.Cause(err) == pgx.ErrNoRows {
				_, err = ctx.Sendf("No term with that ID found.")
				return
			}
			return bot.DB.InternalError(ctx, err)
		}
		exact = true
	} else {
		terms, err := bot.DB.TermName(ctx.RawArgs)
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		} else if err == nil && len(terms) > 0 {
			if len(terms) > 1 {
				return bot.search(ctx)
			}

			exact = true
			term = terms[0]
			goto found
		}

		{
			terms, err := bot.DB.Search(ctx.RawArgs, 1, nil)
			if err != nil {
				return bot.DB.InternalError(ctx, err)
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

	e := bot.DB.TermEmbed(term)

	_, err = m.Embeds(e).Send()
	return
}

func (bot *Bot) termSlash(ctx bcr.Contexter) (err error) {
	query := ctx.GetStringFlag("query")

	if query == "" {
		return ctx.SendEphemeral("You didn't give a term name or ID.")
	}

	var (
		exact bool
		term  *db.Term
	)

	id, err := strconv.Atoi(query)
	if err == nil {
		term, err = bot.DB.GetTerm(id)
		if err != nil {
			if errors.Cause(err) == pgx.ErrNoRows {
				return ctx.SendEphemeral("No term with that ID found.")
			}
			return bot.DB.InternalError(ctx, err)
		}
		exact = true
	} else {
		terms, err := bot.DB.TermName(query)
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		} else if err == nil && len(terms) > 0 {
			if len(terms) > 1 {
				return bot.searchSlash(ctx)
			}

			exact = true
			term = terms[0]
			goto found
		}

		{
			terms, err := bot.DB.Search(query, 1, nil)
			if err != nil {
				return bot.DB.InternalError(ctx, err)
			}
			if len(terms) == 0 {
				return ctx.SendEphemeral("No term found.")
			}

			term = terms[0]
		}
	}

found:
	s := ""

	if !exact {
		s = "I couldn't find a term exactly matching that name, but here's the closest match:"
	}

	e := bot.DB.TermEmbed(term)

	return ctx.SendX(s, e)
}
