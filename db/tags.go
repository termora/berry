package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

// Tags gets all tags from the database
func (db *Db) Tags() (s []string, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &s, "select distinct unnest(tags) as tags from public.terms order by tags")
	return
}

// TagTerms ...
func (db *Db) TagTerms(tag string) (t []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.image_url from public.terms as t, public.categories as c
	where $1 ilike any(t.tags) and t.category = c.id order by t.name, t.id`, tag)
	return
}

// UntaggedTerms ...
func (db *Db) UntaggedTerms() (t []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.image_url from public.terms as t, public.categories as c
	where t.tags = array[]::text[] and t.category = c.id order by t.name, t.id`)
	return
}
