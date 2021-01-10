package db

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

// TermFlag ...
type TermFlag int

// Constants for term flags
const (
	FlagSearchHidden TermFlag = 1 << iota
	FlagRandomHidden
	FlagShowWarning
	FlagListHidden
)

// Term holds info on a single term
type Term struct {
	ID              int       `json:"id"`
	Category        int       `json:"category_id"`
	CategoryName    string    `json:"category"`
	Name            string    `json:"name"`
	Aliases         []string  `json:"aliases"`
	Description     string    `json:"description"`
	Note            string    `json:"note"`
	Source          string    `json:"source"`
	Created         time.Time `json:"created"`
	LastModified    time.Time `json:"last_modified"`
	ContentWarnings string    `json:"content_warnings,omitempty"`

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
func (t *Term) TermEmbed(baseURL string) *discord.Embed {
	var (
		desc = t.Description
		cw   = t.ContentWarnings
		note = t.Note
	)

	if baseURL != "" {
		desc = strings.ReplaceAll(desc, "(##", "("+baseURL)
		note = strings.ReplaceAll(note, "(##", "("+baseURL)
		cw = strings.ReplaceAll(cw, "(##", "("+baseURL)
	}

	defer AddCount()

	fields := make([]discord.EmbedField, 0)
	if len(t.Aliases) != 0 {
		fields = append(fields, discord.EmbedField{
			Name:  "Synonyms",
			Value: strings.Join(t.Aliases, ", "),
		})
	}

	if cw != "" {
		desc = "||" + t.Description + "||"
		fields = append(fields, discord.EmbedField{
			Name:  "Content warning",
			Value: cw,
		})
	}

	if note != "" {
		fields = append(fields, discord.EmbedField{
			Name:  "Note",
			Value: note,
		})
	}

	if t.Warning() {
		fields = append(fields, discord.EmbedField{
			Name:  "Warning",
			Value: "This term is only in this glossary for the sake of completeness. It may be derogatory, exclusionary, or harmful, especially when applied to other people and not as a self-description. Use this term with extreme caution.",
		})
	}

	fields = append(fields, discord.EmbedField{
		Name:  "Source",
		Value: t.Source,
	})

	var u string
	if baseURL != "" {
		u = baseURL + url.PathEscape(strings.ToLower(t.Name))
	}

	e := &discord.Embed{
		Title:       t.Name,
		URL:         u,
		Description: desc,
		Color:       EmbedColour,
		Timestamp:   discord.NewTimestamp(t.Created),
		Fields:      fields,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Category: %v (ID: %v) | Created", t.ID, t.CategoryName, t.Category),
		},
	}

	return e
}
