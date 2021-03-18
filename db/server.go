package db

import "context"

// CreateServerIfNotExists returns true if the server exists
func (db *Db) CreateServerIfNotExists(guildID string) (exists bool, err error) {
	err = db.Pool.QueryRow(context.Background(), "select exists (select from public.servers where id = $1)", guildID).Scan(&exists)
	if err != nil {
		return exists, err
	}
	if !exists {
		commandTag, err := db.Pool.Exec(context.Background(), "insert into public.servers (id, prefixes) values ($1, $2)", guildID, db.Config.Bot.Prefixes)
		if err != nil {
			return exists, err
		}
		if commandTag.RowsAffected() != 1 {
			return exists, ErrorNoRowsAffected
		}
		return exists, err
	}
	return exists, nil
}

// DeleteServer deletes a server's database entry
func (db *Db) DeleteServer(guildID string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "delete from public.servers where id = $1", guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}
