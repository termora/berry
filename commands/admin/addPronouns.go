package admin

import (
	"strings"

	"github.com/starshine-sys/berry/db"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) addPronouns(ctx *bcr.Context) (err error) {
	i := 0
	for _, arg := range strings.Split(ctx.RawArgs, "\n") {
		p := strings.Split(arg, "/")
		if len(p) < 5 {
			_, err = ctx.Sendf("Not enough forms given (argument %v)", arg)
			break
		}

		_, err = c.DB.AddPronoun(db.PronounSet{
			Subjective: p[0],
			Objective:  p[1],
			PossDet:    p[2],
			PossPro:    p[3],
			Reflexive:  p[4],
		})
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
		i++
	}
	_, err = ctx.Sendf("Added %v new pronoun set(s).", i)
	return
}
