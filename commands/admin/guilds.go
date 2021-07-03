package admin

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Admin) cmdGuilds(ctx *bcr.Context) (err error) {
	b := make([]string, 0)

	for _, g := range c.guilds {
		b = append(b, fmt.Sprintf(
			"Name = %v\nID = %v", g.Name, g.ID,
		))
	}
	s := strings.Join(b, "\n\n")

	// if the whole thing fits in a Discord message, send it as that
	// used to be formatted as ini but quotation marks break that
	if len(s) <= 2000 {
		_, err = ctx.Send("", discord.Embed{
			Title:       fmt.Sprintf("Guilds (%v)", len(c.guilds)),
			Description: "```\n" + s + "\n```",
			Color:       db.EmbedColour,
		})
		return err
	}

	// otherwise, compress and upload it
	fn := "guilds.txt"
	buf := new(bytes.Buffer)
	zw := gzip.NewWriter(buf)
	zw.Name = fn
	_, err = zw.Write([]byte(s))
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	err = zw.Close()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	fn += ".gz"

	file := sendpart.File{
		Name:   fn,
		Reader: buf,
	}

	_, err = ctx.State.SendMessageComplex(ctx.Channel.ID, api.SendMessageData{
		Content:         "Here you go!",
		Files:           []sendpart.File{file},
		AllowedMentions: ctx.Router.DefaultMentions,
	})
	return err
}
