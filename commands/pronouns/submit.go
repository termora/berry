package pronouns

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) submit(ctx *bcr.Context) (err error) {
	if c.Config.Bot.Support.PronounChannel == 0 {
		_, err = ctx.Send("We aren't accepting new pronoun submissions through the bot. You might be able to ask in the bot support server.", nil)
		return err
	}

	if _, err = c.submitCooldown.Get(ctx.Author.ID.String()); err == nil {
		_, err = ctx.Send("You can only submit a pronoun set every ten seconds.", nil)
		return
	}

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a pronoun set.", nil)
		return err
	}
	p := strings.Split(ctx.RawArgs, "/")
	if len(p) < 5 {
		_, err = ctx.Send("You didn't give enough forms. Make sure you separate the forms with forward slashes (/).", nil)
		return
	}

	_, err = ctx.NewMessage().Channel(c.Config.Bot.Support.PronounChannel).
		Embed(&discord.Embed{
			Author: &discord.EmbedAuthor{
				Name: fmt.Sprintf("%v#%v (%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
				Icon: ctx.Author.AvatarURL(),
			},
			Color:       db.EmbedColour,
			Title:       "Pronoun submission",
			Description: strings.Join(p[:5], "/"),
			Fields: []discord.EmbedField{{
				Name:  "Submitted by",
				Value: ctx.Author.Mention(),
			}},
			Timestamp: discord.NowTimestamp(),
		}).Send()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.NewMessage().Content(
		fmt.Sprintf("Successfully submitted the pronoun set **%v**.", strings.Join(p[:5], "/")),
	).BlockMentions().Send()
	if err != nil {
		c.Report(ctx, err)
		return err
	}

	c.submitCooldown.SetWithTTL(ctx.Author.ID.String(), true, 10*time.Second)
	return
}
