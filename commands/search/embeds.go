package search

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/termora/berry/db"
)

var emoji = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}

func searchResultEmbed(search string, page, total, totalTerms int, s []*db.Term) discord.Embed {
	var (
		desc   string
		fields []discord.EmbedField
	)

	// only show this if there's more than one page
	if totalTerms > 5 {
		desc = "Use ⬅️ ➡️ to navigate between pages and the numbers to choose a term.\nYou can also type out the number in chat to choose a term."
	}

	// add ellipses to headlines if they're needed
	// this isn't 100% accurate but it's close enough
	for i, t := range s {
		h := t.Headline
		if !strings.HasPrefix(t.Description, h[:10]) && !strings.HasPrefix(t.Headline, "...") {
			h = "..." + h
		}
		if !strings.HasSuffix(t.Description, h[len(h)-10:]) && !strings.HasSuffix(t.Headline, "...") {
			h = h + "..."
		}

		name := t.Name
		if len(t.Aliases) > 0 {
			name += fmt.Sprintf(" (%v)", strings.Join(t.Aliases, ", "))
		}

		fields = append(fields, discord.EmbedField{
			Name:  "​",
			Value: fmt.Sprintf("%v **%v**\n%v\n\n", emoji[i], name, t.Headline),
		})
	}

	return discord.Embed{
		Title:       fmt.Sprintf("Search results for \"%v\"", search),
		Color:       db.EmbedColour,
		Description: desc,
		Fields:      fields,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Results: %v | Page %v/%v", totalTerms, page, total),
		},
	}
}
