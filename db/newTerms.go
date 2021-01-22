package db

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

// TermsSince returns all terms added since the specified date
func (db *Db) TermsSince(d time.Time) (t []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &t, `select t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.content_warnings
	from public.terms as t, public.categories as c
	where t.category = c.id and t.created > $1
	order by name asc`, d)
	return
}
