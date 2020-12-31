package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TermFlag ...
type TermFlag int

// Constants for term flags
const (
	FlagSearchHidden TermFlag = 1 << iota
	FlagRandomHidden
	FlagShowWarning
)

// Term holds info on a single term
type Term struct {
	ID           int       `json:"id"`
	Category     int       `json:"-"`
	CategoryName string    `json:"category"`
	Name         string    `json:"name"`
	Aliases      []string  `json:"aliases"`
	Description  string    `json:"description"`
	Source       string    `json:"source"`
	Created      time.Time `json:"created"`

	Flags TermFlag `json:"flags"`

	// Rank is only populated with db.Search()
	Rank float64 `json:"rank,omitempty"`
	// Headline is only populated with db.Search()
	Headline string `json:"headline,omitempty"`
}

// SearchHidden returns true if the term is hidden from search results
func (t *Term) SearchHidden() bool {
	return t.Flags&FlagSearchHidden == FlagSearchHidden
}

// RandomHidden returns true if the term is hidden from the random command
func (t *Term) RandomHidden() bool {
	return t.Flags&FlagRandomHidden == FlagRandomHidden
}

// Warning returns true if the term has a warning on its term card
func (t *Term) Warning() bool {
	return t.Flags&FlagShowWarning == FlagShowWarning
}

// TermEmbed creates a Discord embed from a term object
func (t *Term) TermEmbed() *discordgo.MessageEmbed {
	if t == nil {
		return nil
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	if len(t.Aliases) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Synonyms",
			Value: strings.Join(t.Aliases, ", "),
		})
	}

	if t.Warning() {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Warning",
			Value: "This term is only in this glossary for the sake of completeness. It may be derogatory, exclusionary, or harmful, especially when applied to other people and not as a self-description. Use this term with extreme caution.",
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "Source",
		Value: t.Source,
	})

	e := &discordgo.MessageEmbed{
		Title:       t.Name,
		Description: t.Description,
		Color:       EmbedColour,
		Timestamp:   t.Created.Format(time.RFC3339),
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %v | Category: %v (ID: %v) | Created", t.ID, t.CategoryName, t.Category),
		},
	}

	return e
}
