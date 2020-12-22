package search

import (
	"fmt"
	"strings"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/termbot/db"
)

func (c *commands) explanation(ctx *crouter.Ctx) (err error) {
	ex, err := c.Db.GetAllExplanations()
	if err != nil {
		return ctx.CommandError(err)
	}
	if ctx.RawArgs != "" {
		for _, e := range ex {
			if strings.ToLower(ctx.RawArgs) == e.Name {
				_, err = ctx.Send(e.Description)
				return err
			}
			for _, alias := range e.Aliases {
				if strings.ToLower(ctx.RawArgs) == alias {
					_, err = ctx.Send(e.Description)
					return err
				}
			}
		}
	}

	var x string
	for _, e := range ex {
		x += fmt.Sprintf("- `%v`\n", e.Name)
	}
	if x == "" {
		x = "No explanations."
	}
	_, err = ctx.Embed("All explanations", x, db.EmbedColour)
	return err
}
