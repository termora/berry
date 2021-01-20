package db

import (
	"context"
	"fmt"

	"github.com/Starshine113/bcr"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/google/uuid"
)

// InternalError sends an error message and logs the error to the database
func (db *Db) InternalError(ctx *bcr.Context, e error) error {
	id := uuid.New()

	_, err := db.Pool.Exec(context.Background(), "insert into public.errors (id, command, user_id, channel, error) values ($1, $2, $3, $4, $5)", id, ctx.Command, ctx.Author.ID, ctx.Channel.ID, e.Error())
	if err != nil {
		// if there's a non-nil error, panic, which should bring us back to the router
		// if the write to the database failed chances are something is *very* wrong anyway
		panic(err)
	}

	_, err = ctx.Send(
		fmt.Sprintf("Error code: ``%v``", bcr.EscapeBackticks(id.String())),
		&discord.Embed{
			Title:       "Internal error occurred",
			Description: "An internal error has occurred. If this issue persists, please contact the bot developer with the error code above.",
			Color:       0xe01a00,

			Footer: &discord.EmbedFooter{
				Text: id.String(),
			},
			Timestamp: discord.NowTimestamp(),
		},
	)
	return err
}
