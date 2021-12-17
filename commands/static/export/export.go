package export

import (
	"time"

	dbpkg "github.com/termora/berry/db"
)

// ExportVersion is the current version
const ExportVersion = 3

// Export is an export of the database
type Export struct {
	Version      int                  `json:"export_version"`
	ExportDate   time.Time            `json:"export_date"`
	Categories   []*dbpkg.Category    `json:"categories"`
	Terms        []*dbpkg.Term        `json:"terms"`
	Tags         []string             `json:"tags"`
	Explanations []*dbpkg.Explanation `json:"explanations,omitempty"`
	Pronouns     []*dbpkg.PronounSet  `json:"pronouns,omitempty"`
}

// New exports the database db
func New(db *dbpkg.DB) (e Export, err error) {
	e = Export{ExportDate: time.Now().UTC(), Version: ExportVersion}

	e.Categories, err = db.GetCategories()
	if err != nil {
		return
	}

	e.Terms, err = db.GetTerms(0)
	if err != nil {
		return
	}

	e.Tags, err = db.Tags()
	if err != nil {
		return
	}

	e.Explanations, err = db.GetAllExplanations()
	if err != nil {
		return
	}

	e.Pronouns, err = db.Pronouns(dbpkg.AlphabeticPronounOrder)
	return
}
