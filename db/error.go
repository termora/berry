package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/common/log"
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
func (db *DB) InternalError(ctx bcr.Contexter, e error) error {
	if db.useSentry {
		return db.sentryError(ctx, e)
	}
	// log to console
	log.Error(e)

	s := "An internal error has occurred. If this issue persists, please contact the bot developer."
	if db.Config != nil {
		if db.Config.Bot.Support.Invite != "" {
			s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v).", db.Config.Bot.Support.Invite)
		}
	}

	_, err := ctx.Send(
		"",
		discord.Embed{
			Title:       "Internal error occurred",
			Description: s,
			Color:       0xE74C3C,
			Timestamp:   discord.NowTimestamp(),
		},
	)
	return err
}

// CaptureError captures an error with additional context
func (db *DB) CaptureError(ctx bcr.Contexter, e error) *sentry.EventID {
	// clone the hub
	hub := db.sentry.Clone()

	// add the user's ID
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: ctx.User().ID.String()})
	})

	guildID := discord.GuildID(0)
	if ctx.GetGuild() != nil {
		guildID = ctx.GetGuild().ID
	}

	cmd := ""
	if v, ok := ctx.(*bcr.Context); ok {
		cmd = v.Command
	} else if v, ok := ctx.(*bcr.SlashContext); ok {
		cmd = v.CommandName
	}

	// add some more info
	hub.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "cmd",
		Data: map[string]interface{}{
			"user":    ctx.User().ID,
			"channel": ctx.GetChannel().ID,
			"guild":   guildID,
			"command": cmd,
		},
		Level:     sentry.LevelError,
		Timestamp: time.Now().UTC(),
	}, nil)

	return hub.CaptureException(e)
}

func (db *DB) sentryError(ctx bcr.Contexter, e error) error {
	log.Error(e)

	// check if it's a problem on our end, to avoid blowing through Sentry's limits
	if !IsOurProblem(e) {
		s := "An internal error has occurred. However, it's unlikely that it's on our end. Please check the input you gave the command again; if you're reasonably sure the error *is* on our end, please contact the bot developer"
		if db.Config.Bot.Support.Invite != "" {
			s = fmt.Sprintf("%v in the [support server](<%v>) with the error code above.", s, db.Config.Bot.Support.Invite)
		} else {
			// hacky as shit but it works :blobsilly:
			s += "."
		}
		_, err := ctx.Send("", discord.Embed{
			Title:       "Internal error occurred",
			Color:       EmbedColour,
			Description: s,
		})
		return err
	}

	id := db.CaptureError(ctx, e)

	s := "An internal error has occurred. If this issue persists, please contact the bot developer with the error code above."
	if db.Config != nil {
		if db.Config.Bot.Support.Invite != "" {
			s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v) with the error code above.", db.Config.Bot.Support.Invite)
		}
	}

	_, err := ctx.Send(
		fmt.Sprintf("Error code: ``%v``", bcr.EscapeBackticks(string(*id))),
		discord.Embed{
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
func (db *DB) Error(id string) (e *Error, err error) {
	e = &Error{}

	ctx, cancel := db.Context()
	defer cancel()

	err = pgxscan.Get(ctx, db.Pool, e, `select
	id, command, user_id, channel, error, time
	from public.errors where id = $1`, id)
	return e, err
}

// SetSentry ...
func (db *DB) SetSentry(hub *sentry.Hub) {
	if hub == nil {
		return
	}
	db.sentry = hub
	db.useSentry = true
}

// IsOurProblem checks if an error is "our problem", as in, should be in the logs and reported to Sentry.
// Will be expanded eventually once we get more insight into what type of errors we get.
func IsOurProblem(e error) bool {
	switch e.(type) {
	case *strconv.NumError:
		// this is because the user inputted an invalid number for string conversion
		// we should handle this in the command itself instead but we're lazy, and this shouldn't come up in normal usage, only with admin commands
		return false
	case *httputil.HTTPError:
		// usually caused by a message being deleted while we're still doing stuff with it (so if someone selects an option in the search results before the bot is done adding reactions)
		return false
	}

	// ignore some specific errors
	switch e {
	case bcr.ErrBotMissingPermissions:
		return false
	case bcr.ErrorNotEnoughArgs, bcr.ErrorTooManyArgs, bcr.ErrInvalidMention, bcr.ErrChannelNotFound, bcr.ErrMemberNotFound, bcr.ErrUserNotFound, bcr.ErrRoleNotFound:
		// we're not sure if these are ever returned, but ignore them anyway
		return false
	}

	return true
}
