package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	cat, terms, err := bot.termCat(strings.Join(ctx.Args, " "))
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	showFullList, _ := ctx.Flags.GetBool("full")
	showAsFile, _ := ctx.Flags.GetBool("file")
	if showFullList {
		return bot.fullList(ctx, terms, showAsFile)
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

func (bot *Bot) termCat(cat string) (s *db.Category, t []*db.Term, err error) {
	if cat != "" {
		id, err := bot.DB.CategoryID(cat)
		if err == nil {
			t, err = bot.DB.GetCategoryTerms(id, search.FlagSearchHidden)
			return bot.DB.CategoryFromID(id), t, err
		}
	}
	t, err = bot.DB.GetTerms(search.FlagSearchHidden)
	return nil, t, err
}

func (bot *Bot) fullList(ctx bcr.Contexter, terms []*db.Term, showAsFile bool) (err error) {
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

		return ctx.SendFiles("", sendpart.File{Name: "list.txt", Reader: strings.NewReader(buf)})
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

	if v, ok := ctx.(*bcr.Context); ok {
		m, _, err := v.PagedEmbedTimeout(embeds, false, time.Hour)
		if err != nil {
			return err
		}

		err = v.State.React(m.ChannelID, m.ID, "❌")
		return err
	} else {
		_, _, err = ctx.ButtonPages(embeds, time.Hour)
		return err
	}
}
