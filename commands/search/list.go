package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search"
)

func (c *commands) list(ctx *bcr.Context) (err error) {
	cat, terms, err := c.termCat(strings.Join(ctx.Args, " "))
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	showFullList, _ := ctx.Flags.GetBool("full")
	showAsFile, _ := ctx.Flags.GetBool("file")
	if showFullList {
		return c.fullList(ctx, terms, showAsFile)
	}

	if showAsFile {
		var buf string

		for _, t := range terms {
			buf += strings.Join(append([]string{t.Name}, t.Aliases...), ", ") + "\n"
		}

		_, err = ctx.NewMessage().AddFile("list.txt", strings.NewReader(strings.TrimSpace(buf))).Send()
		return
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
			t, err = c.DB.GetCategoryTerms(id, search.FlagSearchHidden)
			return c.DB.CategoryFromID(id), t, err
		}
	}
	t, err = c.DB.GetTerms(search.FlagSearchHidden)
	return nil, t, err
}

func (c *commands) fullList(ctx *bcr.Context, terms []*db.Term, showAsFile bool) (err error) {
	if showAsFile {
		var buf string

		buf += fmt.Sprintf("%v terms", len(terms))

		for _, t := range terms {
			buf += fmt.Sprintf(`
-------------------------------
%v`, t.Name)

			if len(t.Aliases) > 0 {
				buf += fmt.Sprintf(", %v", strings.Join(t.Aliases, ", "))
			}
			buf += fmt.Sprintf("\n\n%v\n", t.Description)
		}

		_, err = ctx.NewMessage().AddFile("list.txt", strings.NewReader(buf)).Send()
		return err
	}

	var (
		charCount int
		termCount int

		buf    discord.Embed
		embeds []discord.Embed
		fields []discord.EmbedField
	)

	buf = discord.Embed{
		Title:       fmt.Sprintf("List of terms (%v)", len(terms)),
		Color:       db.EmbedColour,
		Description: "Use ⬅️ ➡️ to flip through pages, and use ❌ to delete this message.",
		Footer:      &discord.EmbedFooter{},
	}

	for _, t := range terms {
		b := fmt.Sprintf("**▶️ %v", t.Name)
		if len(t.Aliases) > 0 {
			b += ", " + strings.Join(t.Aliases, ", ")
		}
		b += "**\n"
		b += t.Description

		if len(b) >= 1010 {
			b = b[:1000] + "..."
		}

		fields = append(fields, discord.EmbedField{Name: "​", Value: b})
	}

	for _, f := range fields {
		charCount += len(f.Value)
		termCount++
		if charCount >= 5000 || termCount > 20 {
			embeds = append(embeds, buf)
			buf = discord.Embed{
				Title:       fmt.Sprintf("List of terms (%v)", len(terms)),
				Color:       db.EmbedColour,
				Description: "Use ⬅️ ➡️ to flip through pages, and use ❌ to delete this message.",
				Footer:      &discord.EmbedFooter{},
			}
			buf.Fields = append(buf.Fields, f)
			charCount = len(f.Value)
			termCount = 1
			continue
		}

		buf.Fields = append(buf.Fields, f)
	}

	embeds = append(embeds, buf)

	for i := range embeds {
		embeds[i].Footer.Text = fmt.Sprintf("Page %v/%v", i+1, len(embeds))
	}

	m, _, err := ctx.PagedEmbedTimeout(embeds, false, time.Hour)
	if err != nil {
		return err
	}

	err = ctx.State.React(m.ChannelID, m.ID, "❌")
	return
}
