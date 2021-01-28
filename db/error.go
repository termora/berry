package db

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/starshine-sys/bcr"
)

// Error ...
type Error struct {
	ID      uuid.UUID
	Command string
	UserID  discord.UserID
	Channel discord.ChannelID
	Error   string
	Time    time.Time
}

// InternalError sends an error message and logs the error to the database
func (db *Db) InternalError(ctx *bcr.Context, e error) error {
	id := uuid.New()

	_, err := db.Pool.Exec(context.Background(), "insert into public.errors (id, command, user_id, channel, error) values ($1, $2, $3, $4, $5)", id, ctx.Command, ctx.Author.ID, ctx.Channel.ID, e.Error())
	if err != nil {
		// if there's a non-nil error, panic, which should bring us back to the router
		// if the write to the database failed chances are something is *very* wrong anyway
		panic(err)
	}

	s := "An internal error has occurred. If this issue persists, please contact the bot developer with the error code above."
	if db.Config != nil {
		if db.Config.Bot.Support.Invite != "" {
			s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v) with the error code above.", db.Config.Bot.Support.Invite)
		}
	}

	_, err = ctx.Send(
		fmt.Sprintf("Error code: ``%v``", bcr.EscapeBackticks(id.String())),
		&discord.Embed{
			Title:       "Internal error occurred",
			Description: s,
			Color:       0xE74C3C,

			Footer: &discord.EmbedFooter{
				Text: id.String(),
			},
			Timestamp: discord.NowTimestamp(),
		},
	)
	return err
}

// Error ...
func (db *Db) Error(id string) (e *Error, err error) {
	e = &Error{}

	err = pgxscan.Get(context.Background(), db.Pool, e, `select
	id, command, user_id, channel, error, time
	from public.errors where id = $1`, id)
	return e, err
}
