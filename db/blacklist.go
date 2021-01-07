package db

import (
	"context"
	"errors"

	"github.com/Starshine113/bcr"
)

// Errors for setting the blacklist
var (
	ErrorAlreadyBlacklisted = errors.New("channel is already blacklisted")
	ErrorNotBlacklisted     = errors.New("channel is not blacklisted")
)

// IsBlacklisted returns true if a channel is blacklisted
func (db *Db) IsBlacklisted(guildID, channelID string) (b bool) {
	db.Pool.QueryRow(context.Background(), "select $1 = any(server.blacklist) from (select * from public.servers where id = $2) as server", channelID, guildID).Scan(&b)
	return b
}

// AddToBlacklist adds the given channelID to the blacklist for guildID
func (db *Db) AddToBlacklist(guildID, channelID string) (err error) {
	err = db.CreateServerIfNotExists(guildID)
	if err != nil {
		return err
	}

	if db.IsBlacklisted(guildID, channelID) {
		return ErrorAlreadyBlacklisted
	}
	commandTag, err := db.Pool.Exec(context.Background(), "update public.servers set blacklist = array_append(blacklist, $1) where id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}

// RemoveFromBlacklist removes the given channelID from the blacklist for guildID
func (db *Db) RemoveFromBlacklist(guildID, channelID string) (err error) {
	err = db.CreateServerIfNotExists(guildID)
	if err != nil {
		return err
	}

	if !db.IsBlacklisted(guildID, channelID) {
		return ErrorNotBlacklisted
	}
	commandTag, err := db.Pool.Exec(context.Background(), "update public.servers set blacklist = array_remove(blacklist, $1) where id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return err
}

// GetBlacklist returns the channel blacklist for guildID
func (db *Db) GetBlacklist(guildID string) (b []string, err error) {
	err = db.CreateServerIfNotExists(guildID)
	if err != nil {
		return b, err
	}

	err = db.Pool.QueryRow(context.Background(), "select blacklist from public.servers where id = $1", guildID).Scan(&b)
	return b, err
}

// CtxInBlacklist is a wrapper around IsBlacklisted for bcr
func (db *Db) CtxInBlacklist(ctx *bcr.Context) bool {
	return db.IsBlacklisted(ctx.Message.GuildID.String(), ctx.Channel.ID.String())
}
