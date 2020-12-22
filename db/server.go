package db

import "context"

// CreateServerIfNotExists ...
func (db *Db) CreateServerIfNotExists(guildID string) (err error) {
	var exists bool
	err = db.Pool.QueryRow(context.Background(), "select exists (select from public.servers where id = $1)", guildID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		commandTag, err := db.Pool.Exec(context.Background(), "insert into public.servers (id) values ($1)", guildID)
		if err != nil {
			return err
		}
		if commandTag.RowsAffected() != 1 {
			return ErrorNoRowsAffected
		}
		return err
	}
	return nil
}
