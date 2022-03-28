package search

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/dustin/go-humanize/english"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *Bot) listTermsSlash(ctx bcr.Contexter) error {
	cat, terms, err := bot.termCat(ctx.GetStringFlag("category"))
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	showFullList := ctx.GetBoolFlag("full")
	showAsFile := ctx.GetBoolFlag("file")
	if showFullList {
		return bot.fullList(ctx, terms, showAsFile)
	}

	if showAsFile {
		var buf string

		for _, t := range terms {
			buf += strings.Join(append([]string{t.Name}, t.Aliases...), ", ") + "\n"
		}

		return ctx.SendFiles("", sendpart.File{Name: "list.txt", Reader: strings.NewReader(buf)})
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

	_, _, err = ctx.ButtonPages(
		PaginateStrings(s, 15, title, "\n"), 15*time.Minute,
	)
	return err
}

func (bot *Bot) listPronounsSlash(ctx bcr.Contexter) error {
	footerTmpl := "Page %v/%v"
	order := db.AlphabeticPronounOrder

	switch ctx.GetStringFlag("sort-by") {
	case "random":
		order = db.RandomPronounOrder
		footerTmpl = "Sorting randomly | Page %v/%v"
	case "alphabetical":
		order = db.AlphabeticPronounOrder
		footerTmpl = "Sorting alphabetically | Page %v/%v"
	case "uses":
		order = db.UsesPronounOrder
		footerTmpl = "Sorting by # of uses | Page %v/%v"
	}

	p, err := bot.DB.Pronouns(order)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
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

	_, _, err = ctx.ButtonPages(e, 15*time.Minute)
	return err
}
