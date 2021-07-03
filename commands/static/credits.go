package static

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Commands) credits(ctx *bcr.Context) (err error) {
	// return if there's no credit fields
	if len(c.Config.Bot.CreditFields) == 0 {
		return nil
	}

	_, err = ctx.Send("", discord.Embed{
		Color:       db.EmbedColour,
		Title:       "Credits",
		Description: fmt.Sprintf("These are the people who helped create %v!", ctx.Bot.Username),
		Fields:      c.Config.Bot.CreditFields,
	})
	return err
}
