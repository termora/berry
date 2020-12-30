package search

import "github.com/Starshine113/crouter"

func (c *commands) random(ctx *crouter.Ctx) (err error) {
	t, err := c.Db.RandomTerm()
	if err != nil {
		return ctx.CommandError(err)
	}
	_, err = ctx.Send(t.TermEmbed())
	return
}
