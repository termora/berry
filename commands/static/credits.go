package static

import (
	"fmt"

	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
)

func (c *Commands) credits(ctx *crouter.Ctx) (err error) {
	// return if there's no credit fields
	if len(c.config.Bot.CreditFields) == 0 {
		return nil
	}

	fs := make([]*discordgo.MessageEmbedField, 0)

	for _, f := range c.config.Bot.CreditFields {
		fs = append(fs, &discordgo.MessageEmbedField{
			Name:  f.Name,
			Value: f.Value,
		})
	}

	_, err = ctx.Send(&discordgo.MessageEmbed{
		Color:       db.EmbedColour,
		Title:       "Credits",
		Description: fmt.Sprintf("These are the people who helped create %v!", ctx.BotUser.Username),
		Fields:      fs,
	})
	return err
}
