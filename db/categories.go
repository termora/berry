package db

// AddCategory ...
func (db *Db) AddCategory(name string) (id int, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Adding category %v", name)

	err = db.QueryRow(ctx, "insert into public.categories (name) values ($1) returning id", name).Scan(&id)

	Debug("Added category %v", id)
	return
}
