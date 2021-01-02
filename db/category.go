package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

// Category is a single category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CategoryID gets the ID from a category name
func (db *Db) CategoryID(s string) (id int, err error) {
	err = db.Pool.QueryRow(context.Background(), "select id from public.categories where lower(name) = lower($1)", s).Scan(&id)
	return
}

// GetCategories ...
func (db *Db) GetCategories() (c []*Category, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &c, `select id, name
	from public.categories`)
	return c, err
}
