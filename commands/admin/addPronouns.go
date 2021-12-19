package admin

import (
	"strings"

	"github.com/termora/berry/db"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) addPronouns(ctx *bcr.Context) (err error) {
	i, skipped := 0, 0
	for _, arg := range strings.Split(ctx.RawArgs, "\n") {
		p := strings.Split(arg, "/")
		if len(p) < 5 {
			skipped++
			continue
		}

		_, err = bot.DB.AddPronoun(db.PronounSet{
			Subjective: p[0],
			Objective:  p[1],
			PossDet:    p[2],
			PossPro:    p[3],
			Reflexive:  p[4],
		})
		if err != nil {
			skipped++
			continue
		}
		i++
	}
	_, err = ctx.Sendf("Added %v new pronoun set(s) (skipped %v)", i, skipped)
	return
}
