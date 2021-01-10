package admin

import (
	"strings"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/misc"
)

func (c *commands) aio(ctx *bcr.Context) (err error) {
	if ctx.CheckRequiredArgs(5); err != nil {
		_, err = ctx.Send("Too few or too many arguments supplied.", nil)
		return err
	}

	name := ctx.Args[0]
	catName := ctx.Args[1]
	description := ctx.Args[2]
	var aliases []string
	if ctx.Args[3] == "none" {
		aliases = []string{}
	} else {
		aliases = strings.Split(ctx.Args[3], ",")
	}
	source := ctx.Args[4]

	category, err := c.db.CategoryID(catName)
	if err != nil {
		_, err = ctx.Send("Could not find that category, cancelled.", nil)
		return
	}
	if category == 0 {
		return
	}

	t := &db.Term{
		Name:        name,
		Category:    category,
		Description: description,
		Aliases:     aliases,
		Source:      source,
	}

	msg, err := ctx.Send("Term finished. React with ✅ to finish adding it, or with ❌ to cancel. Preview:", t.TermEmbed(""))
	if err != nil {
		return err
	}

	ctx.AddYesNoHandler(*msg, ctx.Author.ID, func(ctx *bcr.Context) {
		t, err := c.db.AddTerm(t)
		if err != nil {
			_, err = ctx.Send(misc.InternalError, nil)
			return
		}
		ctx.Sendf("Added term with ID %v.", t.ID)
	}, func(ctx *bcr.Context) {
		ctx.Send("Cancelled.", nil)
	})

	return nil
}
