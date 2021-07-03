package search

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) tags(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		t, err := c.DB.Tags()
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}

		_, err = ctx.PagedEmbed(PaginateStrings(t, 15, "Tags", "\n"), false)
		return err
	}

	terms, err := c.DB.TagTerms(ctx.RawArgs)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("I couldn't find any terms with that tag.")
			return
		}
		return c.DB.InternalError(ctx, err)
	}

	if len(terms) == 0 {
		_, err = ctx.Send("I couldn't find any terms with that tag.")
		return
	}

	var s []string
	for _, t := range terms {
		s = append(s, t.Name)
	}

	_, err = ctx.PagedEmbed(
		PaginateStrings(s, 15, fmt.Sprintf("Terms tagged ``%v``", bcr.EscapeBackticks(ctx.RawArgs)), "\n"), false,
	)
	return
}

// PaginateStrings paginates the given slice of strings as embeds.
func PaginateStrings(slice []string, itemsPerPage int, title, join string) (embeds []discord.Embed) {
	var (
		b     []string
		count int
		buf   [][]string
	)

	for _, s := range slice {
		if count >= itemsPerPage {
			buf = append(buf, b)
			count = 0
			b = nil
		}

		b = append(b, s)
		count++
	}

	buf = append(buf, b)

	for i, s := range buf {
		embeds = append(embeds, discord.Embed{
			Title:       fmt.Sprintf("%v (%v)", title, len(slice)),
			Description: strings.Join(s, join),
			Color:       db.EmbedColour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+1, len(buf)),
			},
		})
	}

	return
}
