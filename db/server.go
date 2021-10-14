package db

// CreateServerIfNotExists returns true if the server exists
func (db *Db) CreateServerIfNotExists(guildID string) (exists bool, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	err = db.QueryRow(ctx, "select exists (select from public.servers where id = $1)", guildID).Scan(&exists)
	if err != nil {
		return exists, err
	}
	if !exists {
		ctx, cancel := db.Context()
		defer cancel()

		Debug("Creating server %v", guildID)

		commandTag, err := db.Exec(ctx, "insert into public.servers (id, prefixes) values ($1, $2)", guildID, db.Config.Bot.Prefixes)
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
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Deleting server %v", guildID)

	commandTag, err := db.Exec(ctx, "delete from public.servers where id = $1", guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}
