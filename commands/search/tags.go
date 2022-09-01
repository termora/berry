package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *Bot) tags(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		t, err := bot.DB.Tags()
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		}

		_, err = ctx.PagedEmbed(PaginateStrings(t, 15, "Tags", "\n"), false)
		return err
	}

	terms, err := bot.DB.TagTerms(ctx.RawArgs)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("I couldn't find any terms with that tag.")
			return
		}
		return bot.DB.InternalError(ctx, err)
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

func (bot *Bot) tagsSlash(ctx bcr.Contexter) (err error) {
	if ctx.GetStringFlag("tag") == "" {
		t, err := bot.DB.Tags()
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		}

		for i := range t {
			t[i] = t[i] + "\n"
		}

		_, _, err = ctx.ButtonPages(bcr.StringPaginator("Tags", db.EmbedColour, t, 15), 15*time.Minute)
		return err
	}

	tag := ctx.GetStringFlag("tag")
	terms, err := bot.DB.TagTerms(tag)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return ctx.SendX("I couldn't find any terms with that tag.")
		}
		return bot.DB.InternalError(ctx, err)
	}

	if len(terms) == 0 {
		_, err = ctx.Send("I couldn't find any terms with that tag.")
		return
	}

	var s []string
	for _, t := range terms {
		s = append(s, t.Name+"\n")
	}

	_, _, err = ctx.ButtonPages(
		bcr.StringPaginator("Terms tagged \""+tag+"\"", db.EmbedColour, s, 15), 15*time.Minute,
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
