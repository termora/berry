package server

import (
	"fmt"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/berry/db"
)

func (c *commands) blacklist(ctx *crouter.Ctx) (err error) {
	b, err := c.db.GetBlacklist(ctx.Message.GuildID)
	if err != nil {
		return ctx.CommandError(err)
	}
	var x string
	for _, c := range b {
		x += fmt.Sprintf("<#%v>\n", c)
	}
	if len(b) == 0 {
		x = "No channels are blacklisted."
	}
	_, err = ctx.Embed("Channel blacklist", x, db.EmbedColour)
	return err
}
