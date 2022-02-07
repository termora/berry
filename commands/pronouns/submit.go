package pronouns

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
)

func (bot *Bot) submit(ctx *bcr.Context) (err error) {
	if bot.Config.Bot.PronounChannel == 0 {
		_, err = ctx.Send("We aren't accepting new pronoun submissions through the bot. You might be able to ask in the bot support server.")
		return err
	}

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You didn't give a pronoun set.")
		return err
	}
	p := strings.Split(ctx.RawArgs, "/")
	if len(p) < 5 {
		_, err = ctx.Send("You didn't give enough forms. Make sure you separate the forms with forward slashes (/).")
		return
	}
	if len(p) > 5 {
		_, err = ctx.Send("You gave too many forms. Make sure you have five forms, separated with forward slashes.")
		return
	}

	// normalize pronouns
	for i := range p {
		p[i] = strings.ToLower(strings.TrimSpace(p[i]))
	}

	_, err = bot.DB.GetPronoun(p...)
	if err == nil {
		_, err = ctx.Send("That pronoun set already exists!")
		return
	}

	con, cancel := bot.DB.Context()
	defer cancel()

	found := false
	err = bot.DB.QueryRow(con, "select exists(select * from pronoun_msgs where subjective = $1 and objective = $2 and poss_det = $3 and poss_pro = $4 and reflexive = $5)", p[0], p[1], p[2], p[3], p[4]).Scan(&found)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	if found {
		_, err = ctx.Send("That pronoun set has already been submitted!")
		return
	}

	msg, err := ctx.NewMessage().Channel(bot.Config.Bot.PronounChannel).
		Embeds(discord.Embed{
			Author: &discord.EmbedAuthor{
				Name: fmt.Sprintf("%v#%v (%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
				Icon: ctx.Author.AvatarURL(),
			},
			Color:       db.EmbedColour,
			Title:       "Pronoun submission",
			Description: strings.Join(p, "/"),
			Fields: []discord.EmbedField{{
				Name:  "Submitted by",
				Value: ctx.Author.Mention(),
			}},
			Timestamp: discord.NowTimestamp(),
		}).Send()
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	con, cancel = bot.DB.Context()
	defer cancel()

	_, err = bot.DB.Exec(con, "insert into pronoun_msgs (message_id, subjective, objective, poss_det, poss_pro, reflexive) values ($1, $2, $3, $4, $5, $6)", msg.ID, p[0], p[1], p[2], p[3], p[4])
	if err == nil {
		// if the error's non-nil, the message was still sent
		// so don't just return immediately
		ctx.State.React(msg.ChannelID, msg.ID, "✅")
	} else {
		log.Errorf("Error adding submission message %v to database: %v", msg.ID, err)
	}

	_, err = ctx.NewMessage().Content(
		fmt.Sprintf("Successfully submitted the pronoun set **%v**.", strings.Join(p[:5], "/")),
	).BlockMentions().Send()
	if err != nil {
		bot.Report(ctx, err)
		return err
	}

	return
}
