package search

import (
	"fmt"
	"strings"

	"github.com/Starshine113/berry/db"
	"github.com/bwmarrin/discordgo"
)

var emoji = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}

func searchResultEmbed(search string, page, total int, s []*db.Term) *discordgo.MessageEmbed {
	var desc string
	for i, t := range s {
		h := t.Headline
		if !strings.HasPrefix(t.Description, h[:10]) {
			h = "..." + h
		}
		if !strings.HasSuffix(t.Description, h[len(h)-10:]) {
			h = h + "..."
		}
		name := t.Name
		if len(t.Aliases) > 0 {
			name += fmt.Sprintf(" (%v)", strings.Join(t.Aliases, ", "))
		}
		desc += fmt.Sprintf("%v **%v** (rank: %v)\n%v\n\n", emoji[i], name, t.Rank, h)
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Search results for \"%v\"", search),
		Description: desc,
		Color:       db.EmbedColour,
		Fields: []*discordgo.MessageEmbedField{{
			Name:   "Usage",
			Value:  "Use ⬅️ ➡️ to navigate between pages, the numbers to choose a term, and ❌ to delete this message.",
			Inline: false,
		}},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %v/%v", page, total),
		},
	}
}
