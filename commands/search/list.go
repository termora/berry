package search

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) list(ctx *bcr.Context) (err error) {
	cat, terms, err := c.termCat(ctx.RawArgs)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	s := make([]string, 0)
	for _, t := range terms {
		if len(t.Aliases) > 0 {
			s = append(s, fmt.Sprintf("`%v`: %v, %v", t.ID, t.Name, strings.Join(t.Aliases, ", ")))
			continue
		}
		s = append(s, fmt.Sprintf("`%v`: %v", t.ID, t.Name))
	}

	// create pages of slices
	termSlices := make([][]string, 0)
	// 15 terms each
	for i := 0; i < len(s); i += 15 {
		end := i + 15

		if end > len(s) {
			end = len(s)
		}

		termSlices = append(termSlices, s[i:end])
	}

	// create the embeds and send them
	embeds := make([]discord.Embed, 0)

	title := fmt.Sprintf("List of terms (%v)", len(terms))
	footer := ""
	if cat != nil {
		title = fmt.Sprintf("List of %v terms (%v)", cat.Name, len(terms))
		footer = fmt.Sprintf("Category: %v (ID: %v) |", cat.Name, cat.ID)
	}
	for i, s := range termSlices {
		embeds = append(embeds, discord.Embed{
			Title:       title,
			Description: strings.Join(s, "\n"),
			Color:       db.EmbedColour,

			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("%v Page %v/%v", footer, i+1, len(termSlices)),
			},
		})
	}

	_, err = ctx.PagedEmbed(embeds, false)
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
