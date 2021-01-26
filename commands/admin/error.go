package admin

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (c *Admin) error(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You didn't give an error ID.", nil)
		return err
	}

	e, err := c.db.Error(ctx.RawArgs)
	if err != nil {
		c.sugar.Errorf("Error when retrieving error with ID %v: %v", ctx.RawArgs, err)
		_, err = ctx.Send("Error with that ID not found, or another error occurred.", nil)
		return err
	}

	_, err = ctx.Send("", &discord.Embed{
		Title:       e.ID.String(),
		Description: "```" + e.Error + "```",
		Fields: []discord.EmbedField{{
			Name:  "Context",
			Value: fmt.Sprintf("- **Command:** %v\n- **User:** %v\n- **Channel:** %v", e.Command, e.UserID, e.Channel),
		}},
		Footer: &discord.EmbedFooter{
			Text: e.ID.String(),
		},
		Color:     0xE74C3C,
		Timestamp: discord.NewTimestamp(e.Time),
	})
	return err
}
