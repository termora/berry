package admin

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
)

// shows a list of all guilds
func (c *commands) guilds(ctx *crouter.Ctx) (err error) {
	var b bytes.Buffer
	for _, g := range ctx.Session.State.Guilds {
		b.WriteString(fmt.Sprintf("%v (%v)\n", g.ID, g.Name))
	}

	u, err := ctx.Session.UserChannelCreate(ctx.Author.ID)
	if err != nil {
		c.sugar.Error("Error creating DM channel:", err)
		_, err = ctx.Send("Error creating DM channel.")
		return
	}

	if len(b.String()) < 2000 {
		_, err = ctx.Session.ChannelMessageSend(u.ID, fmt.Sprintf("```Guilds (%v)\n=============\n%v```", len(ctx.Session.State.Guilds), b.String()))
		return
	}

	file := discordgo.File{
		Name:   fmt.Sprintf("guilds-%v-%v.txt", ctx.BotUser.Username, time.Now().Format("2006-01-02-15-04-05")),
		Reader: &b,
	}

	_, err = ctx.Session.ChannelMessageSendComplex(u.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Done! List of %v guilds:", len(ctx.Session.State.Guilds)),
		Files:   []*discordgo.File{&file},
	})
	return err
}
