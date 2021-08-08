package typesense

import (
	"fmt"

	"github.com/termora/berry/db/search"
	"github.com/typesense/typesense-go/typesense/api"
)

// SyncTerms synchronizes the given terms with the Typesense server.
func (c *Client) SyncTerms(terms []*search.Term) error {
	// delete the existing collection, if any
	existing := c.ts.Collection("terms")
	if existing != nil {
		existing.Delete()
	}

	schema := &api.CollectionSchema{
		Name: "terms",
		Fields: []api.Field{
			{
				Name: "names",
				Type: "string[]",
			},
			{
				Name: "description",
				Type: "string",
			},
			{
				Name: "source",
				Type: "string",
			},
			{
				Name:  "tags",
				Type:  "string[]",
				Facet: true,
			},
			{
				Name:  "category",
				Type:  "string",
				Facet: true,
			},
		},
	}

	_, err := c.ts.Collections().Create(schema)
	if err != nil {
		return err
	}

	docs := []interface{}{}

	for _, t := range terms {
		docs = append(docs, tsTerm{
			ID:          t.ID,
			Category:    t.Category,
			Names:       append([]string{t.Name}, t.Aliases...),
			Description: t.Description,
			Source:      t.Source,
			Tags:        t.Tags,
		})
	}

	_, err = c.ts.Collection("terms").Documents().Import(docs, &api.ImportDocumentsParams{Action: "upsert"})
	return err
}

type tsTerm struct {
	ID          int      `json:"id,string"`
	Category    int      `json:"category,string"`
	Names       []string `json:"names"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Tags        []string `json:"tags"`
}

// SyncTerm upserts a single term.
func (c *Client) SyncTerm(t *search.Term) error {
	doc := tsTerm{
		ID:          t.ID,
		Category:    t.Category,
		Names:       append([]string{t.Name}, t.Aliases...),
		Description: t.Description,
		Source:      t.Source,
		Tags:        t.Tags,
	}

	_, err := c.ts.Collection("terms").Documents().Upsert(doc)
	return err
}

// SyncDelete deletes a single term.
func (c *Client) SyncDelete(id int) error {
	doc := c.ts.Collection("terms").Document(fmt.Sprint(id))
	if doc != nil {
		_, err := doc.Delete()
		return err
	}
	return nil
}

func boolPointer(b bool) *bool { return &b }
