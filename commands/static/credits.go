package static

import (
	"fmt"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
)

func (c *Commands) credits(ctx *bcr.Context) (err error) {
	// return if there's no credit fields
	if len(c.config.Bot.CreditFields) == 0 {
		return nil
	}

	fs := make([]discord.EmbedField, 0)

	for _, f := range c.config.Bot.CreditFields {
		fs = append(fs, discord.EmbedField{
			Name:  f.Name,
			Value: f.Value,
		})
	}

	_, err = ctx.Send("", &discord.Embed{
		Color:       db.EmbedColour,
		Title:       "Credits",
		Description: fmt.Sprintf("These are the people who helped create %v!", ctx.Bot.Username),
		Fields:      fs,
	})
	return err
}
