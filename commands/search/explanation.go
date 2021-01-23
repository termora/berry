package search

import (
	"fmt"
	"strings"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
)

func (c *commands) explanation(ctx *bcr.Context) (err error) {
	ex, err := c.Db.GetAllExplanations()
	if err != nil {
		return c.Db.InternalError(ctx, err)
	}
	if ctx.RawArgs != "" {
		for _, e := range ex {
			if strings.ToLower(ctx.RawArgs) == e.Name {
				_, err = ctx.Send(e.Description, nil)
				return err
			}
			for _, alias := range e.Aliases {
				if strings.ToLower(ctx.RawArgs) == alias {
					_, err = ctx.Send(e.Description, nil)
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
	_, err = ctx.Send("", &discord.Embed{
		Title:       "All explanations",
		Description: x,
		Color:       db.EmbedColour,
	})
	return err
}
