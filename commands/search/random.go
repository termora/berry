package search

import (
	"github.com/starshine-sys/bcr"
)

func (c *commands) random(ctx *bcr.Context) (err error) {
	// grab a random term
	t, err := c.DB.RandomTerm()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	// send the random term
	_, err = ctx.Send("", t.TermEmbed(c.Config.Bot.TermBaseURL))
	return
}
