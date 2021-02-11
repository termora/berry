package search

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
)

func (c *commands) list(ctx *bcr.Context) (err error) {
	terms, err := c.termCat(ctx.RawArgs)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	s := make([]string, 0)
	for _, t := range terms {
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

	for i, s := range termSlices {
		embeds = append(embeds, discord.Embed{
			Title:       fmt.Sprintf("List of terms (%v)", len(terms)),
			Description: strings.Join(s, "\n"),
			Color:       db.EmbedColour,

			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(termSlices)),
			},
		})
	}

	_, err = ctx.PagedEmbed(embeds, false)
	return err
}

func (c *commands) termCat(cat string) (t []*db.Term, err error) {
	if cat != "" {
		cat, err := c.DB.CategoryID(cat)
		if err == nil {
			return c.DB.GetCategoryTerms(cat, db.FlagSearchHidden)
		}
	}
	return c.DB.GetTerms(db.FlagSearchHidden)
}
