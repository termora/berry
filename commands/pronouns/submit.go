package pronouns

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
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

	_, err = ctx.Router.Session.SendMessageComplex(c.Config.Bot.Support.PronounChannel, api.SendMessageData{
		Content:         fmt.Sprintf("%v#%v (<@%v>, %v) submitted a new pronoun set: **%v**.", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID, ctx.Author.ID, strings.Join(p[:5], "/")),
		AllowedMentions: &api.AllowedMentions{Parse: nil},
	})
	if err != nil {
		id := c.Report(err)
		var e *discord.Embed
		if id != nil {
			e = &discord.Embed{
				Title:       "Error code",
				Description: "`" + string(*id) + "`",
				Fields: []discord.EmbedField{{
					Name:  "Internal error occurred",
					Value: "To report this error, please join the bot support server and give the bot developer the ID above.",
				}},
				Color: db.EmbedColour,
			}
		}
		_, err = ctx.Send("There was an error sending the submission. Please try again, and if this happens again, please submit it in the support server instead.", e)
		if err != nil {
			c.Report(err)
			return err
		}
	}

	_, err = ctx.Send("", &discord.Embed{
		Description: fmt.Sprintf("Successfully submitted the pronoun set **%v**.", strings.Join(p[:5], "/")),
		Color:       db.EmbedColour,
	})
	if err != nil {
		c.Report(err)
		return err
	}

	c.submitCooldown.SetWithTTL(ctx.Author.ID.String(), true, 5*time.Minute)
	return
}
