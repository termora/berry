package pronouns

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) list(ctx *bcr.Context) (err error) {
	p, err := c.DB.Pronouns()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	var (
		b     strings.Builder
		s     = make([]string, 0)
		count int
	)

	for _, p := range p {
		if count >= 20 {
			s = append(s, b.String())
			b.Reset()
			count = 0
		}
		b.WriteString(fmt.Sprintf("`%v`: %s\n", p.ID, p))
		count++
	}
	s = append(s, b.String())

	e := make([]discord.Embed, 0)
	for i, page := range s {
		e = append(e, discord.Embed{
			Title:       fmt.Sprintf("List of pronouns (%v)", len(p)),
			Description: page,
			Color:       db.EmbedColour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(s)),
			},
		})
	}

	_, err = ctx.PagedEmbed(e, false)
	return
}
