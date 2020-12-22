package db

import "context"

// AddCategory ...
func (db *Db) AddCategory(name string) (id int, err error) {
	err = db.Pool.QueryRow(context.Background(), "insert into public.categories (name) values ($1) returning id", name).Scan(&id)
	return
}
