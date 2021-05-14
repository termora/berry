package db

import "errors"

// Searcher is an interface containing all methods needed to search for and return terms.
type Searcher interface {
	// RefreshTerms refreshes the search index with terms
	RefreshTerms(terms []*Term) (err error)

	// Search searches the database for terms
	Search(input string, limit int, ignore []string) (terms []*Term, err error)

	// SearchCat searches for terms from a single category
	SearchCat(input string, cat, limit int, showHidden bool, ignore []string) (terms []*Term, err error)
}

// noopSearcher is a no-op default Searcher that always errors.
type noopSearcher int

var errNoopSearcher = errors.New("db: no-op searcher called")

var _ Searcher = (*noopSearcher)(nil)

func (noopSearcher) RefreshTerms(terms []*Term) (err error) { return errNoopSearcher }
func (noopSearcher) Search(input string, limit int, ignore []string) (terms []*Term, err error) {
	return nil, errNoopSearcher
}
func (noopSearcher) SearchCat(input string,
	cat, limit int,
	showHidden bool, ignore []string) (terms []*Term, err error) {
	return nil, errNoopSearcher
}

// NewNoopSearcher returns a new Searcher that is no-op.
// No idea why you'd ever want this but it's there I guess.
func NewNoopSearcher() Searcher { return noopSearcher(0) }
