package search

import (
	"github.com/Starshine113/bcr"
)

func (c *commands) random(ctx *bcr.Context) (err error) {
	t, err := c.Db.RandomTerm()
	if err != nil {
		return c.Db.InternalError(ctx, err)
	}
	_, err = ctx.Send("", t.TermEmbed(c.conf.Bot.TermBaseURL))
	return
}
