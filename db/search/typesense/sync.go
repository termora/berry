package typesense

import (
	"fmt"

	"github.com/termora/berry/db/search"
	"github.com/termora/tsclient"
)

// SyncTerms synchronizes the given terms with the Typesense server.
func (c *Client) SyncTerms(terms []*search.Term) error {
	_, err := c.ts.DeleteCollection("terms")
	if err != nil && err != tsclient.ErrNotFound {
		return err
	}

	_, err = c.ts.CreateCollection("terms", "", []tsclient.CreateFieldData{
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
	})
	if err != nil {
		return err
	}

	docs := []tsTerm{}

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

	_, err = c.ts.Import("terms", "upsert", docs)
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

	return c.ts.Upsert("terms", doc, nil)
}

// SyncDelete deletes a single term.
func (c *Client) SyncDelete(id int) error {
	return c.ts.DeleteDocument("terms", fmt.Sprint(id), nil)
}

func boolPointer(b bool) *bool { return &b }
