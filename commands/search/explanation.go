package search

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) explanation(ctx *bcr.Context) (err error) {
	ex, err := c.DB.GetAllExplanations()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	// create a new message object
	m := ctx.NewMessage()
	if ctx.Message.Reference != nil {
		m = m.Reference(ctx.Message.Reference.MessageID)

		var o bool = false

		m = m.AllowedMentions(&api.AllowedMentions{
			Parse: []api.AllowedMentionType{api.AllowUserMention},

			RepliedUser: option.Bool(&o),
		})

		if len(ctx.Message.Mentions) > 0 {
			o = true
			m.Data.AllowedMentions.RepliedUser = option.Bool(&o)
		}
	}

	// just cycle through all of these, it's fine (probably)
	if ctx.RawArgs != "" {
		for _, e := range ex {
			if strings.EqualFold(ctx.RawArgs, e.Name) {
				_, err = m.Content(e.Description).Send()
				return err
			}
			for _, alias := range e.Aliases {
				if strings.EqualFold(ctx.RawArgs, alias) {
					_, err = m.Content(e.Description).Send()
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

	_, err = m.Embeds(discord.Embed{
		Title:       "All explanations",
		Description: x,
		Color:       db.EmbedColour,
	}).Send()
	return err
}
