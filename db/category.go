package db

import (
	"github.com/georgysavva/scany/pgxscan"
)

// Category is a single category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CategoryID gets the ID from a category name
func (db *Db) CategoryID(s string) (id int, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting category with name %v", s)

	err = db.QueryRow(ctx, "select id from public.categories where lower(name) = lower($1)", s).Scan(&id)
	return
}

// GetCategories ...
func (db *Db) GetCategories() (c []*Category, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting categories")

	err = pgxscan.Select(ctx, db.Pool, &c, `select id, name
	from public.categories`)
	return c, err
}

// CategoryFromID ...
func (db *Db) CategoryFromID(id int) (c *Category) {
	c = &Category{}

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting category with ID %v", id)

	pgxscan.Get(ctx, db.Pool, c, `select id, name
	from public.categories where id = $1`, id)
	return c
}
