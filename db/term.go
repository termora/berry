package db

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
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
	Note            string    `json:"note,omitempty"`
	Source          string    `json:"source"`
	Created         time.Time `json:"created"`
	LastModified    time.Time `json:"last_modified"`
	Tags            []string  `json:"-"`
	DisplayTags     []string  `json:"tags,omitempty"`
	ContentWarnings string    `json:"content_warnings,omitempty"`
	ImageURL        string    `json:"image_url,omitempty"`

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
	defer AddCount()

	Debug("Creating term embed for %v", t.ID)

	e := &discord.Embed{
		Title:     t.Name,
		Color:     EmbedColour,
		Timestamp: discord.NewTimestamp(t.Created),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Category: %v (ID: %v) | Created", t.ID, t.CategoryName, t.Category),
		},
	}

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

	if cw != "" {
		desc = "||" + desc + "||"

		if len(desc) < 1024 {
			e.Description = fmt.Sprintf("**Content warning: %v**", cw)
		} else {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "​",
				Value: fmt.Sprintf("**Content warning: %v**", cw),
			})
		}
	}

	if len(desc) < 1024 && cw != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Description",
			Value: desc,
		})
	} else {
		e.Description = desc
	}

	if len(t.Aliases) != 0 {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Synonyms",
			Value: strings.Join(t.Aliases, ", "),
		})
	}

	if note != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Note",
			Value: note,
		})
	}

	if t.Warning() {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Warning",
			Value: "This term is only in this glossary for the sake of completeness. It may be derogatory, exclusionary, or harmful, especially when applied to other people and not as a self-description. Use this term with extreme caution.",
		})
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name:  "Source",
		Value: t.Source,
	})

	if len(t.DisplayTags) > 0 {
		var b strings.Builder
		for i, tag := range t.DisplayTags {
			if b.Len() >= 500 {
				b.WriteString(fmt.Sprintf("\nToo many to list (showing %v/%v)", i, len(t.DisplayTags)))
				break
			}
			b.WriteString(tag)
			if i != len(t.DisplayTags)-1 {
				b.WriteString(", ")
			}
		}

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Tag(s)",
			Value: b.String(),
		})
	}

	if baseURL != "" {
		e.URL = baseURL + url.PathEscape(strings.ToLower(t.Name))
	}

	if t.ImageURL != "" {
		e.Image = &discord.EmbedImage{
			URL: t.ImageURL,
		}
	}

	return e
}
