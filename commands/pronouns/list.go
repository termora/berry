package pronouns

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize/english"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) list(ctx *bcr.Context) (err error) {
	footerTmpl := "Page %v/%v"
	order := db.AlphabeticPronounOrder
	if flag, _ := ctx.Flags.GetBool("random"); flag {
		order = db.RandomPronounOrder
		footerTmpl = "Sorting randomly | Page %v/%v"
	} else if flag, _ := ctx.Flags.GetBool("alphabetical"); flag {
		order = db.AlphabeticPronounOrder
		footerTmpl = "Sorting alphabetically | Page %v/%v"
	} else if flag, _ := ctx.Flags.GetBool("by-uses"); flag {
		order = db.UsesPronounOrder
		footerTmpl = "Sorting by # of uses | Page %v/%v"
	}

	p, err := c.DB.Pronouns(order)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	if order == db.RandomPronounOrder {
		rand.Shuffle(len(p), func(i, j int) {
			p[i], p[j] = p[j], p[i]
		})
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
		b.WriteString(p.String())
		if order == db.UsesPronounOrder {
			b.WriteString(" (" + english.Plural(int(p.Uses), "use", "uses") + ")")
		}
		b.WriteRune('\n')
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
				Text: fmt.Sprintf(footerTmpl, i+1, len(s)),
			},
		})
	}

	_, err = ctx.PagedEmbed(e, false)
	return
}
