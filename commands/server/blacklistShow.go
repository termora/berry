package server

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"

	"github.com/starshine-sys/berry/db"
)

func (c *commands) blacklist(ctx *bcr.Context) (err error) {
	b, err := c.db.GetBlacklist(ctx.Message.GuildID.String())
	if err != nil {
		_, err = ctx.Sendf(":x: Internal error occurred: %v", err)
		return err
	}
	var x string
	// append all channel IDs (as mentions) to x
	for _, c := range b {
		x += fmt.Sprintf("<#%v>\n", c)
	}
	if len(b) == 0 {
		x = "No channels are blacklisted."
	}
	_, err = ctx.Send("", &discord.Embed{
		Title:       "Channel blacklist",
		Description: x,
		Color:       db.EmbedColour,
	})
	return err
}
