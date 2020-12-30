package misc

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// PagedEmbedBuilder ...
type PagedEmbedBuilder struct {
	AuthorIconURL string
	AuthorName    string
	Title         string
	Color         int
	embeds        []*discordgo.MessageEmbed
}

// NewEmbedBuilder ...
func NewEmbedBuilder(title, authorName, authorIconURL string, color int) *PagedEmbedBuilder {
	return &PagedEmbedBuilder{
		AuthorIconURL: authorIconURL,
		AuthorName:    authorName,
		Title:         title,
		Color:         color,
		embeds:        make([]*discordgo.MessageEmbed, 0),
	}
}

// Add ...
func (p *PagedEmbedBuilder) Add(title, desc string, fields []*discordgo.MessageEmbedField) *PagedEmbedBuilder {
	if title == "" {
		title = p.Title
	}
	p.embeds = append(p.embeds, &discordgo.MessageEmbed{
		Title: title,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: p.AuthorIconURL,
			Name:    p.AuthorName,
		},
		Color:       p.Color,
		Description: desc,
		Fields:      fields,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	})

	return p
}

// Build finalizes the embed page numbers and returns the slice
func (p *PagedEmbedBuilder) Build() []*discordgo.MessageEmbed {
	for i, e := range p.embeds {
		e.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %v/%v", i+1, len(p.embeds)),
		}
	}

	return p.embeds
}
