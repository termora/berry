package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) explanation(ctx bcr.Contexter) (err error) {
	name := ctx.GetStringFlag("explanation")
	if name == "" {
		if v, ok := ctx.(*bcr.Context); ok {
			name = v.RawArgs
		}
	}

	var text string
	err = c.DB.Pool.QueryRow(context.Background(), "select description from explanations where $1 ilike any(aliases) or $1 ilike name", name).Scan(&text)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {

		}
		return c.DB.InternalError(ctx, err)
	}

	v, ok := ctx.(*bcr.Context)
	if ok {
		var reference *discord.MessageReference
		if v.Message.Reference != nil {
			reference = &discord.MessageReference{MessageID: v.Message.Reference.MessageID}
		}
		_, err = ctx.Session().SendMessageComplex(ctx.GetChannel().ID, api.SendMessageData{
			Content:   text,
			Reference: reference,
		})
		return
	}

	return ctx.SendX(text)
}

func (c *commands) listAllExplanations(ctx bcr.Contexter) (err error) {
	ex, err := c.DB.GetAllExplanations()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	s := []string{}
	for _, e := range ex {
		s = append(s, fmt.Sprintf("- %v\n", strings.Join(append([]string{e.Name}, e.Aliases...), ", ")))
	}

	_, _, err = ctx.ButtonPages(
		bcr.StringPaginator(fmt.Sprintf("All explanations (%v)", len(ex)), db.EmbedColour, s, 10),
		15*time.Minute,
	)
	return
}
