package search

import (
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) random(ctx bcr.Contexter) (err error) {
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
		b, err := bot.randomCategory(ctx, catName, ignore)
		if b || err != nil {
			return err
		}
	}

	// grab a random term
	t, err := bot.DB.RandomTerm(ignore)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return ctx.SendEphemeral("No terms found! Are you sure you're not excluding every possible term?")
		}
		return bot.DB.InternalError(ctx, err)
	}

	// send the random term
	_, err = ctx.Send("", bot.DB.TermEmbed(t))
	return
}

func (bot *Bot) randomCategory(ctx bcr.Contexter, catName string, ignore []string) (b bool, err error) {
	cat, err := bot.DB.CategoryID(catName)
	if err != nil {
		// dont bother to check if its a category not found error or not, just return nil
		return false, nil
	}

	t, err := bot.DB.RandomTermCategory(cat, ignore)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			err = ctx.SendEphemeral("No terms found! Are you sure you're not excluding every possible term?")
			return true, err
		}
		return true, bot.DB.InternalError(ctx, err)
	}

	err = ctx.SendX("", bot.DB.TermEmbed(t))
	return true, err
}
