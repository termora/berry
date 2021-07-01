package db

// AddCategory ...
func (db *Db) AddCategory(name string) (id int, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	err = db.Pool.QueryRow(ctx, "insert into public.categories (name) values ($1) returning id", name).Scan(&id)
	return
}
