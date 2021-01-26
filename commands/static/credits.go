package static

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
)

func (c *Commands) credits(ctx *bcr.Context) (err error) {
	// return if there's no credit fields
	if len(c.Config.Bot.CreditFields) == 0 {
		return nil
	}

	fs := make([]discord.EmbedField, 0)

	for _, f := range c.Config.Bot.CreditFields {
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
