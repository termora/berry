package pronouns

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) submit(ctx *bcr.Context) (err error) {
	if c.Config.Bot.Support.PronounChannel == 0 {
		_, err = ctx.Send("We aren't accepting new pronoun submissions through the bot. You might be able to ask in the bot support server.", nil)
		return err
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
	if len(p) > 5 {
		_, err = ctx.Send("You gave too many forms. Make sure you have five forms, separated with forward slashes.", nil)
		return
	}

	_, err = c.DB.GetPronoun(strings.Split(ctx.RawArgs, "/")...)
	if err == nil {
		_, err = ctx.Send("That pronoun set already exists!", nil)
		return
	}

	msg, err := ctx.NewMessage().Channel(c.Config.Bot.Support.PronounChannel).
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

	_, err = c.DB.Pool.Exec(context.Background(), "insert into pronoun_msgs (message_id, subjective, objective, poss_det, poss_pro, reflexive) values ($1, $2, $3, $4, $5, $6)", msg.ID, p[0], p[1], p[2], p[3], p[4])
	if err == nil {
		// if the error's non-nil, the message was still sent
		// so don't just return immediately
		ctx.State.React(msg.ChannelID, msg.ID, "✅")
	} else {
		c.Sugar.Errorf("Error adding submission message %v to database: %v", msg.ID, err)
	}

	_, err = ctx.NewMessage().Content(
		fmt.Sprintf("Successfully submitted the pronoun set **%v**.", strings.Join(p[:5], "/")),
	).BlockMentions().Send()
	if err != nil {
		c.Report(ctx, err)
		return err
	}

	return
}
