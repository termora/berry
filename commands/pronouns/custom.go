package pronouns

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *commands) custom(ctx *bcr.Context) (err error) {
	set := ctx.Args
	if len(set) != 5 {
		set = strings.Split(ctx.RawArgs, "/")
		if len(set) != 5 {
			_, err = ctx.Send("You gave either too few or too many forms, please give exactly 5.", nil)
			return
		}
	}

	if tmplCount == 0 {
		_, err = ctx.Send("There are no examples available for pronouns! If you think this is in error, please join the bot support server and ask there.", nil)
		return err
	}

	var (
		b strings.Builder
		e = make([]discord.Embed, 0)
	)

	use := &db.PronounSet{
		Subjective: set[0],
		Objective:  set[1],
		PossDet:    set[2],
		PossPro:    set[3],
		Reflexive:  set[4],
	}

	e = append(e, discord.Embed{
		Title:       fmt.Sprintf("%v/%v pronouns", set[0], set[1]),
		Description: fmt.Sprintf("**%s**\n\nTo see these pronouns in action, use the arrow reactions on this message!", use),
		Color:       db.EmbedColour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Page 1/%v", tmplCount+1),
		},
	})

	for i := 0; i < tmplCount; i++ {
		err = templates.ExecuteTemplate(&b, strconv.Itoa(i), use)
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		}
		e = append(e, discord.Embed{
			Title:       fmt.Sprintf("%v/%v pronouns", use.Subjective, use.Objective),
			Description: b.String(),
			Color:       db.EmbedColour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Page %v/%v", i+2, tmplCount+1),
			},
		})
		b.Reset()
	}

	_, err = ctx.PagedEmbed(e, false)
	return
}
