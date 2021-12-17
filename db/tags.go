package db

import (
	"github.com/georgysavva/scany/pgxscan"
)

// Tags gets all tags from the database
func (db *DB) Tags() (s []string, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting all tags")

	err = db.QueryRow(ctx, "select array(select display from tags order by tags)").Scan(&s)
	return
}

// TagTerms ...
func (db *DB) TagTerms(tag string) (t []*Term, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting terms with tag %v", tag)

	err = pgxscan.Select(ctx, db.Pool, &t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.image_url from public.terms as t, public.categories as c
	where $1 ilike any(t.tags) and t.category = c.id order by t.name, t.id`, tag)
	return
}

// UntaggedTerms ...
func (db *DB) UntaggedTerms() (t []*Term, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting untagged terms")

	err = pgxscan.Select(ctx, db.Pool, &t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.image_url from public.terms as t, public.categories as c
	where t.tags = array[]::text[] and t.category = c.id order by t.name, t.id`)
	return
}
