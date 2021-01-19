package search

import (
	"time"

	"github.com/Starshine113/bcr"
)

func (c *commands) initExplanations(r *bcr.Router) {
	explanations, err := c.Db.GetCmdExplanations()
	if err != nil {
		c.Sugar.Error("Error getting explanations:", err)
		return
	}

	for _, e := range explanations {
		e := e
		r.AddCommand(&bcr.Command{
			Name:    e.Name,
			Aliases: e.Aliases,

			Cooldown: 1 * time.Second,
			Command: func(ctx *bcr.Context) (err error) {
				_, err = ctx.Send(e.Description, nil)
				return err
			},
		})
	}
}
