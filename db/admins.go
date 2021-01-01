package db

import (
	"context"
)

// AddAdmin adds an admin to the database
func (db *Db) AddAdmin(id string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "insert into public.admins (user_id) values ($1)", id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}

// GetAdmins gets the current admins as a slice of strings
func (db *Db) GetAdmins() (admins []string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select array(select user_id from public.admins)").Scan(&admins)
	return
}
