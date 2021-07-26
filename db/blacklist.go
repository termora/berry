package db

import (
	"errors"

	"github.com/starshine-sys/bcr"
)

// Errors for setting the blacklist
var (
	ErrorAlreadyBlacklisted = errors.New("channel is already blacklisted")
	ErrorNotBlacklisted     = errors.New("channel is not blacklisted")
)

// IsBlacklisted returns true if a channel is blacklisted
func (db *Db) IsBlacklisted(guildID, channelID string) (b bool) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Checking if channel %v is blacklisted", channelID)

	db.Pool.QueryRow(ctx, "select $1 = any(server.blacklist) from (select * from public.servers where id = $2) as server", channelID, guildID).Scan(&b)
	return b
}

// AddToBlacklist adds the given channelID to the blacklist for guildID
func (db *Db) AddToBlacklist(guildID string, channelIDs []string) (err error) {
	for _, channelID := range channelIDs {
		if db.IsBlacklisted(guildID, channelID) {
			return ErrorAlreadyBlacklisted
		}
	}

	ctx, cancel := db.Context()
	defer cancel()

	commandTag, err := db.Pool.Exec(ctx, "update public.servers set blacklist = array_cat(blacklist, $1) where id = $2", channelIDs, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	Debug("Added %v to blacklist", channelIDs)

	return err
}

// RemoveFromBlacklist removes the given channelID from the blacklist for guildID
func (db *Db) RemoveFromBlacklist(guildID, channelID string) (err error) {
	if !db.IsBlacklisted(guildID, channelID) {
		return ErrorNotBlacklisted
	}

	ctx, cancel := db.Context()
	defer cancel()

	commandTag, err := db.Pool.Exec(ctx, "update public.servers set blacklist = array_remove(blacklist, $1) where id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	Debug("Removed %v from blacklist", channelID)

	return err
}

// GetBlacklist returns the channel blacklist for guildID
func (db *Db) GetBlacklist(guildID string) (b []string, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting blacklist for %v", guildID)

	err = db.Pool.QueryRow(ctx, "select blacklist from public.servers where id = $1", guildID).Scan(&b)
	return b, err
}

// CtxInBlacklist is a wrapper around IsBlacklisted for bcr
func (db *Db) CtxInBlacklist(ctx *bcr.Context) bool {
	if ctx.Guild == nil {
		return false
	}

	if db.IsBlacklisted(ctx.Message.GuildID.String(), ctx.Channel.ID.String()) {
		return true
	}

	if !ctx.Thread() {
		return false
	}

	return db.IsBlacklisted(ctx.Message.GuildID.String(), ctx.ParentChannel.ID.String())
}
