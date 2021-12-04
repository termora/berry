package typesense

import (
	"context"
	"strconv"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/termora/berry/db/search"
	"github.com/termora/tsclient"
	"github.com/termora/tsclient/utils/jsonutil"
)

// SearchCat searches a specific category for a term.
func (c *Client) SearchCat(input string, cat, limit int, ignore []string) (terms []*search.Term, err error) {
	var filterBy string
	if cat != 0 {
		filterBy = "category:" + strconv.Itoa(cat)
	}

	c.Debug("Searching for \"%v\"", input)

	resp, err := c.ts.Search("terms", tsclient.SearchData{
		NoPreSegmentedQuery:     true,
		Query:                   input,
		QueryBy:                 []string{"names", "description", "source"},
		QueryByWeights:          []int{2, 1, 1},
		Prefix:                  []bool{true, false, false},
		SortBy:                  []string{"_text_match:desc"},
		NumTypos:                jsonutil.IntPointer(0),
		PerPage:                 limit,
		HighlightStartTag:       jsonutil.StringPointer("**"),
		HighlightEndTag:         jsonutil.StringPointer("**"),
		HighlightAffixNumTokens: 10,
		FilterBy:                filterBy,
	})
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get a single connection for all requests below
	conn, err := c.pg.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	for i, hit := range resp.Hits {
		var doc tsTerm
		err = hit.UnmarshalTo(&doc)
		if err != nil {
			c.Debug("Error getting term index %v: %v", i, err)
			continue
		}

		t, err := c.getTerm(ctx, conn, doc.ID)
		if err != nil {
			c.Debug("Error getting term ID %v: %v", doc.ID, err)
			return nil, err
		}

		// check tags
		for _, tag := range t.Tags {
			for _, ignore := range ignore {
				if tag == ignore {
					continue
				}
			}
		}

		// extract headline
		for _, hl := range hit.Highlights {
			if hl.Field != "description" || hl.Snippet == "" {
				continue
			}
			t.Headline = hl.Snippet
			break
		}

		if t.Headline == "" {
			t.Headline = t.Description
			if len(t.Description) > 103 {
				t.Headline = t.Description[:100] + "..."
			}
		}

		terms = append(terms, t)
	}

	return terms, nil
}

// Search searches the database for terms
func (c *Client) Search(input string, limit int, ignore []string) (terms []*search.Term, err error) {
	return c.SearchCat(input, 0, limit, ignore)
}

// getTerm gets a term by ID, as the output returned from Typesense isn't 100% complete (and also a string -> interface map, ew)
func (c *Client) getTerm(ctx context.Context, conn *pgxpool.Conn, id int) (t *search.Term, err error) {
	t = &search.Term{}

	err = pgxscan.Get(ctx, conn, t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c where t.id = $1 and t.category = c.id`, id)
	return t, err
}
