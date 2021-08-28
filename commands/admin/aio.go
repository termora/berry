package admin

import (
	"fmt"
	"strings"

	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/admin/auditlog"
	"github.com/termora/berry/db"
)

func (c *Admin) aio(ctx *bcr.Context) (err error) {
	// this command requires 5 arguments exactly
	if ctx.CheckRequiredArgs(5); err != nil {
		_, err = ctx.Send("Too few or too many arguments supplied.")
		return err
	}

	// 0: name
	// 1: category
	// 2: description
	// 3: aliases
	// 4: source
	name := ctx.Args[0]
	catName := ctx.Args[1]
	description := ctx.Args[2]

	var aliases []string
	if ctx.Args[3] == "none" {
		aliases = []string{}
	} else {
		aliases = strings.Split(ctx.Args[3], ",")
	}
	for i := range aliases {
		aliases[i] = strings.TrimSpace(aliases[i])
	}

	source := ctx.Args[4]

	category, err := c.DB.CategoryID(catName)
	if err != nil {
		_, err = ctx.Send("Could not find that category, cancelled.")
		return
	}
	if category == 0 {
		return
	}

	// create and add the term object
	t := &db.Term{
		Name:        name,
		Category:    category,
		Description: description,
		Aliases:     aliases,
		Source:      source,
		Tags:        []string{c.DB.CategoryFromID(category).Name},
	}

	t, err = c.DB.AddTerm(t)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = c.AuditLog.SendLog(t.ID, auditlog.TermEntry, auditlog.CreateAction, nil, t, ctx.Author.ID, nil)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send(fmt.Sprintf("Added term with ID %v.", t.ID), c.DB.TermEmbed(t))
	if err != nil {
		c.Report(ctx, err)
	}
	return err
}
