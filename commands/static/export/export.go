package export

import (
	"time"

	"github.com/starshine-sys/berry/db"
)

// ExportVersion is the current version
const ExportVersion = 3

// Export is an export of the database
type Export struct {
	Version      int               `json:"export_version"`
	ExportDate   time.Time         `json:"export_date"`
	Categories   []*db.Category    `json:"categories"`
	Terms        []*db.Term        `json:"terms"`
	Tags         []string          `json:"tags"`
	Explanations []*db.Explanation `json:"explanations,omitempty"`
	Pronouns     []*db.PronounSet  `json:"pronouns,omitempty"`
}

// New exports the database db
func New(db *db.Db) (e Export, err error) {
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

	e.Pronouns, err = db.Pronouns()
	return
}
