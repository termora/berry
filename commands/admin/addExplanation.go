package admin

import (
	"strings"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/berry/db"
)

func (c *commands) addExplanation(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}
	e := &db.Explanation{}

	content := strings.Split(ctx.RawArgs, "\n")
	if len(content) < 2 {
		_, err = ctx.Sendf("Not enough arguments provided.")
		return err
	}
	names := strings.Split(content[0], " ")
	e.Name = names[0]
	if len(names) > 1 {
		e.Aliases = names[1:]
	}
	e.Description = strings.Join(content[1:], "\n")
	e, err = c.db.AddExplanation(e)
	if err != nil {
		return ctx.CommandError(err)
	}
	_, err = ctx.Sendf("Added explanation with ID %v.\nName: `%v`\nAliases: `%v`", e.ID, e.Name, strings.Join(e.Aliases, ", "))
	return err
}
