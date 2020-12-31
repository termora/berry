package admin

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
)

// shows a list of all guilds
func (c *commands) guilds(ctx *crouter.Ctx) (err error) {
	var b bytes.Buffer
	for _, g := range ctx.Session.State.Guilds {
		b.WriteString(fmt.Sprintf("%v (%v)\n", g.ID, g.Name))
	}

	if len(b.String()) < 2000 {
		_, err = ctx.Embed("Guilds", "```"+b.String()+"```", db.EmbedColour)
		return
	}

	file := discordgo.File{
		Name:   fmt.Sprintf("guilds-%v-%v.txt", ctx.BotUser.Username, time.Now().Format("2006-01-02-15-04-05")),
		Reader: &b,
	}

	_, err = ctx.Send(&discordgo.MessageSend{
		Content: "Done! List of guilds:",
		Files:   []*discordgo.File{&file},
	})
	return err
}
