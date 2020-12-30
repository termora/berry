package search

import (
	"fmt"
	"strings"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/termbot/db"
	"github.com/Starshine113/termbot/misc"
)

func (c *commands) list(ctx *crouter.Ctx) (err error) {
	terms, err := c.Db.GetTerms(db.FlagSearchHidden)
	if err != nil {
		return ctx.CommandError(err)
	}
	s := make([]string, 0)
	for _, t := range terms {
		s = append(s, fmt.Sprintf("`%v`: %v", t.ID, t.Name))
	}

	termSlices := make([][]string, 0)

	for i := 0; i < len(s); i += 10 {
		end := i + 5

		if end > len(s) {
			end = len(s)
		}

		termSlices = append(termSlices, s[i:end])
	}

	b := misc.NewEmbedBuilder("List of terms", "", "", db.EmbedColour)
	for _, s := range termSlices {
		b.Add("", strings.Join(s, "\n"), nil)
	}

	_, err = ctx.PagedEmbed(b.Build())
	return err
}
