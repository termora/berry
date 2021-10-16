package typesense

import "github.com/termora/tsclient"

const autocompleteLimit = 25

// Autocomplete ...
func (c *Client) Autocomplete(input string) (terms []string, err error) {
	c.Debug("Invoking autocomplete for \"%v\"", input)

	resp, err := c.ts.Search("terms", tsclient.SearchData{
		Query:            input,
		QueryBy:          []string{"names"},
		SortBy:           []string{"_text_match:desc"},
		PerPage:          autocompleteLimit,
		SnippetThreshold: 1000,
	})
	if err != nil {
		return
	}

	for i, hit := range resp.Hits {
		var doc tsTerm
		err = hit.UnmarshalTo(&doc)
		if err != nil {
			c.Debug("Error getting term index %v: %v", i, err)
			continue
		}
		terms = append(terms, doc.Names[0])
	}
	return terms, nil
}
