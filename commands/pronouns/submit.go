package pronouns

import (
	"fmt"
	"strings"
	"time"

	"github.com/starshine-sys/bcr"
)

func (c *commands) submit(ctx *bcr.Context) (err error) {
	if c.Config.Bot.Support.PronounChannel == 0 {
		_, err = ctx.Send("We aren't accepting new pronoun submissions through the bot. You might be able to ask in the bot support server.", nil)
		return err
	}

	if _, err = c.submitCooldown.Get(ctx.Author.ID.String()); err == nil {
		_, err = ctx.Send("You can only submit a pronoun set every five minutes, please try again later.\nIf you're submitting a lot of pronouns at once, feel free to join the bot support server and ask there!", nil)
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
		Content(
			fmt.Sprintf(
				"%v#%v (<@%v>, %v) submitted a new pronoun set: **%v**.",
				ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID,
				ctx.Author.ID, strings.Join(p[:5], "/"),
			)).BlockMentions().Send()
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

	c.submitCooldown.SetWithTTL(ctx.Author.ID.String(), true, 5*time.Minute)
	return
}
