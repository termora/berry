package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/getsentry/sentry-go"
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
	if db.useSentry {
		return db.sentryError(ctx, e)
	}
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

func (db *Db) sentryError(ctx *bcr.Context, e error) error {
	// check if it's a problem on our end, to avoid blowing through Sentry's limits, and to clean up error logs
	if !IsOurProblem(e) {
		s := "An internal error has occurred. However, it's unlikely that it's on our end. Please check the input you gave the command again; if you're reasonably sure the error *is* on our end, please contact the bot developer"
		if db.Config.Bot.Support.Invite != "" {
			s = fmt.Sprintf("%v in the [support server](<%v>) with the error code above.", s, db.Config.Bot.Support.Invite)
		} else {
			// hacky as shit but it works :blobsilly:
			s += "."
		}
		_, err := ctx.Send("", &discord.Embed{
			Title:       "Internal error occurred",
			Color:       EmbedColour,
			Description: s,
			Fields: []discord.EmbedField{{
				Name:  "Usage",
				Value: fmt.Sprintf("```%v %v```", ctx.Cmd.Name, ctx.Cmd.Usage),
			}},
		})
		return err
	}

	db.Sugar.Error(e)
	id := db.sentry.CaptureException(e)
	if id == nil {
		return nil
	}

	s := "An internal error has occurred. If this issue persists, please contact the bot developer with the error code above."
	if db.Config != nil {
		if db.Config.Bot.Support.Invite != "" {
			s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v) with the error code above.", db.Config.Bot.Support.Invite)
		}
	}

	_, err := ctx.Send(
		fmt.Sprintf("Error code: ``%v``", bcr.EscapeBackticks(string(*id))),
		&discord.Embed{
			Title:       "Internal error occurred",
			Description: s,
			Color:       0xE74C3C,

			Footer: &discord.EmbedFooter{
				Text: string(*id),
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

// SetSentry ...
func (db *Db) SetSentry(hub *sentry.Hub) {
	if hub == nil {
		return
	}
	db.sentry = hub
	db.useSentry = true
}

// IsOurProblem checks if an error is "our problem", as in, should be in the logs and reported to Sentry.
// Will be expanded eventually once we get more insight into what type of errors we get.
func IsOurProblem(e error) bool {
	// first, check types
	switch e.(type) {
	case *strconv.NumError:
		// this is because the user inputted an invalid number for string conversion
		// we should handle this in the command itself instead but we're lazy, and this shouldn't come up in normal usage, only with admin commands
		return false
	case *httputil.HTTPError:
		// 404 error, so just return false
		// usually caused by a message being deleted while we're still doing stuff with it (so if someone selects an option in the search results before the bot is done adding reactions)
		if e.(*httputil.HTTPError).Code == 404 {
			return false
		}
	}

	return true
}
