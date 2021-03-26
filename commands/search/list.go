package search

import (
	"fmt"
	"strings"

	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) list(ctx *bcr.Context) (err error) {
	cat, terms, err := c.termCat(ctx.RawArgs)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	var s []string
	for _, t := range terms {
		s = append(s, strings.Join(
			append([]string{t.Name}, t.Aliases...), ", ",
		))
	}

	title := fmt.Sprintf("List of terms")
	if cat != nil {
		title = fmt.Sprintf("List of %v terms", cat.Name)
	}

	_, err = ctx.PagedEmbed(
		PaginateStrings(s, 15, title, "\n"), false,
	)
	return err
}

func (c *commands) termCat(cat string) (s *db.Category, t []*db.Term, err error) {
	if cat != "" {
		id, err := c.DB.CategoryID(cat)
		if err == nil {
			t, err = c.DB.GetCategoryTerms(id, db.FlagSearchHidden)
			return c.DB.CategoryFromID(id), t, err
		}
	}
	t, err = c.DB.GetTerms(db.FlagSearchHidden)
	return nil, t, err
}
