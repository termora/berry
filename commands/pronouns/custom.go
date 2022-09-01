package pronouns

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *Bot) custom(ctx bcr.Contexter) (err error) {
	set := strings.Split(ctx.GetStringFlag("set"), "/")
	if v, ok := ctx.(*bcr.Context); ok {
		set = v.Args
		if len(set) != 5 {
			set = strings.Split(v.RawArgs, "/")
		}
	}

	if len(set) != 5 {
		return ctx.SendX("You gave either too few or too many forms, please give exactly 5.")
	}

	if tmplCount == 0 {
		return ctx.SendX("There are no examples available for pronouns! If you think this is in error, please join the bot support server and ask there.")
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

	_, _, err = ctx.ButtonPages(e, 15*time.Minute)
	return
}
