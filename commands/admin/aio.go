package admin

import (
	"fmt"
	"strings"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
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

	t, err = c.db.AddTerm(t)
	if err != nil {
		return c.db.InternalError(ctx, err)
	}
	_, err = ctx.Send(fmt.Sprintf("Added term with ID %v.", t.ID), t.TermEmbed(c.config.Bot.TermBaseURL))
	return err
}
