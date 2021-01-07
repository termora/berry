package search

import (
	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/misc"
)

func (c *commands) random(ctx *bcr.Context) (err error) {
	t, err := c.Db.RandomTerm()
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}
	_, err = ctx.Send("", t.TermEmbed(c.conf.Bot.TermBaseURL))
	return
}
