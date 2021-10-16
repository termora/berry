// Package search contains types for searching for terms.
package search

import "time"

// Searcher is an interface for searching the term database.
type Searcher interface {
	Search(input string, limit int, ignore []string) (terms []*Term, err error)
	SearchCat(input string, cat, limit int, ignore []string) (terms []*Term, err error)

	// Return terms starting with the input
	Autocomplete(input string) (terms []string, err error)

	// Synchronize all terms with the search store
	SyncTerms(terms []*Term) (err error)
	// Synchronize a single term after updates
	SyncTerm(term *Term) (err error)
	// Delete a term from the search store
	SyncDelete(id int) (err error)
}

// TermFlag ...
type TermFlag int

// Constants for term flags
const (
	FlagSearchHidden TermFlag = 1 << iota
	FlagRandomHidden
	FlagShowWarning
	FlagListHidden
	FlagDisputed
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

// Disputed returns true if the term is disputed
func (t *Term) Disputed() bool {
	return t.Flags&FlagDisputed == FlagDisputed
}
