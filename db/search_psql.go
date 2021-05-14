package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ Searcher = (*psqlSearcher)(nil)

type psqlSearcher struct {
	*pgxpool.Pool
}

// NewPsqlSearcher creates a new Searcher with a Postgres backend
func NewPsqlSearcher(pool *pgxpool.Pool) Searcher {
	return &psqlSearcher{Pool: pool}
}

// RefreshTerms always returns nil as it is not needed for the Postgres backend
func (psqlSearcher) RefreshTerms(terms []*Term) (err error) { return nil }

// Search searches the database for terms
func (db *psqlSearcher) Search(input string, limit int, ignore []string) (terms []*Term, err error) {
	if limit == 0 {
		limit = 50
	}
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 8) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0
	and not $4 && tags
	order by rank desc
	limit $2`, input, limit, FlagSearchHidden, ignore)
	return terms, err
}

// SearchCat searches for terms from a single category
func (db *psqlSearcher) SearchCat(input string, cat, limit int, showHidden bool, ignore []string) (terms []*Term, err error) {
	if limit == 0 {
		limit = 50
	}

	flags := FlagSearchHidden
	if showHidden {
		flags = 0
	}

	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 8) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0 and t.category = $4
	and not $5 && tags
	order by rank desc
	limit $2`, input, limit, flags, cat, ignore)
	return terms, err
}
