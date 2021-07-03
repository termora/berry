package admin

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Admin) delTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.Send("No term ID provided.")
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	t, err := c.DB.GetTerm(id)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	m, err := ctx.Send("Are you sure you want to delete this term? React with ✅ to delete it, or with ❌ to cancel.", t.TermEmbed(""))
	if err != nil {
		return err
	}

	// confirm deleting the term
	if yes, timeout := ctx.YesNoHandler(*m, ctx.Author.ID); !yes || timeout {
		ctx.Send("Cancelled.")
		return
	}

	err = c.DB.RemoveTerm(id)
	if err != nil {
		c.Sugar.Error("Error removing term:", err)
		c.DB.InternalError(ctx, err)
		return
	}
	_, err = ctx.Send("✅ Term deleted.")
	if err != nil {
		c.Sugar.Error("Error sending message:", err)
	}

	// if logging terms is enabled, log this
	if c.WebhookClient != nil {
		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Content: "​",

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title:     "Term deleted",
					Color:     db.EmbedColour,
					Timestamp: discord.NowTimestamp(),
				},
				t.TermEmbed(""),
			},
		})
	}
	return nil
}
