package admin

import (
	"strings"

	"github.com/starshine-sys/berry/db"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) addPronouns(ctx *bcr.Context) (err error) {
	p := strings.Split(ctx.RawArgs, "/")
	if len(p) < 5 {
		_, err = ctx.Send("Not enough forms given.", nil)
	}

	id, err := c.db.AddPronoun(db.PronounSet{
		Subjective: p[0],
		Objective:  p[1],
		PossDet:    p[2],
		PossPro:    p[3],
		Reflexive:  p[4],
	})

	if err != nil {
		if err == db.ErrNoForms {
			_, err = ctx.Send("Not enough forms given, some were empty.", nil)
			return err
		}
		return c.db.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Added new pronoun set with ID %v.", id)
	return
}
