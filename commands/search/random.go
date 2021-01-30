package search

import (
	"github.com/starshine-sys/bcr"
)

func (c *commands) random(ctx *bcr.Context) (err error) {
	// grab a random term
	t, err := c.Db.RandomTerm()
	if err != nil {
		return c.Db.InternalError(ctx, err)
	}

	// send the random term
	_, err = ctx.Send("", t.TermEmbed(c.conf.Bot.TermBaseURL))
	return
}
