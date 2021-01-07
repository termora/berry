package misc

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

// PagedEmbedBuilder ...
type PagedEmbedBuilder struct {
	AuthorIconURL string
	AuthorName    string
	Title         string
	Color         int
	embeds        []discord.Embed
}

// NewEmbedBuilder ...
func NewEmbedBuilder(title, authorName, authorIconURL string, color int) *PagedEmbedBuilder {
	return &PagedEmbedBuilder{
		AuthorIconURL: authorIconURL,
		AuthorName:    authorName,
		Title:         title,
		Color:         color,
		embeds:        make([]discord.Embed, 0),
	}
}

// Add ...
func (p *PagedEmbedBuilder) Add(title, desc string, fields []discord.EmbedField) *PagedEmbedBuilder {
	if title == "" {
		title = p.Title
	}
	p.embeds = append(p.embeds, discord.Embed{
		Title: title,
		Author: &discord.EmbedAuthor{
			Icon: p.AuthorIconURL,
			Name: p.AuthorName,
		},
		Color:       discord.Color(p.Color),
		Description: desc,
		Fields:      fields,
		Timestamp:   discord.NewTimestamp(time.Now()),
	})

	return p
}

// Build finalizes the embed page numbers and returns the slice
func (p *PagedEmbedBuilder) Build() []discord.Embed {
	for i, e := range p.embeds {
		e.Footer = &discord.EmbedFooter{
			Text: fmt.Sprintf("Page %v/%v", i+1, len(p.embeds)),
		}
	}

	return p.embeds
}
