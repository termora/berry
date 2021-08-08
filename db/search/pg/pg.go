// Package pg implements the search.Searcher interface with the existing PostgreSQL database.
package pg

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/termora/berry/db/search"
)

// New returns a Searcher based on the given pgxpool.Pool
func New(pool *pgxpool.Pool, debugFunc func(string, ...interface{})) search.Searcher {
	if debugFunc == nil {
		debugFunc = func(string, ...interface{}) {
			return
		}
	}
	return &pg{pool, debugFunc}
}

var _ search.Searcher = (*pg)(nil)

type pg struct {
	*pgxpool.Pool
	Debug func(template string, args ...interface{})
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// Search searches the database for terms
func (db *pg) Search(input string, limit int, ignore []string) (terms []*search.Term, err error) {
	if limit == 0 {
		limit = 50
	}

	db.Debug("Searching for terms `%v`, limit %v, ignoring `%v`", input, limit, ignore)

	ctx, cancel := getContext()
	defer cancel()

	err = pgxscan.Select(ctx, db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 8) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0
	and not $4 && tags
	order by rank desc
	limit $2`, input, limit, search.FlagSearchHidden, ignore)
	return terms, err
}

// SearchCat searches for terms from a single category
func (db *pg) SearchCat(input string, cat, limit int, ignore []string) (terms []*search.Term, err error) {
	if limit == 0 {
		limit = 50
	}

	db.Debug("Searching for terms `%v` in category %v, limit %v, ignoring `%v`", input, cat, limit, ignore)

	ctx, cancel := getContext()
	defer cancel()

	err = pgxscan.Select(ctx, db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 8) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0 and t.category = $4
	and not $5 && tags
	order by rank desc
	limit $2`, input, limit, search.FlagSearchHidden, cat, ignore)
	return terms, err
}

// SyncTerms is no-op in the postgres backend.
func (*pg) SyncTerms([]*search.Term) error {
	return nil
}

// SyncTerm is no-op in the postgres backend.
func (*pg) SyncTerm(*search.Term) error {
	return nil
}

// SyncDelete is no-op in the postgres backend.
func (*pg) SyncDelete(int) error {
	return nil
}
