package db

import "context"

// CategoryID gets the ID from a category name
func (db *Db) CategoryID(s string) (id int, err error) {
	err = db.Pool.QueryRow(context.Background(), "select id from public.categories where lower(name) = lower($1)", s).Scan(&id)
	return
}
