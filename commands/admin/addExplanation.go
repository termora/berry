package admin

import (
	"strings"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/misc"
	"github.com/diamondburned/arikawa/v2/bot/extras/shellwords"
)

func (c *commands) addExplanation(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("Not enough arguments provided.", nil)
		return err
	}
	e := &db.Explanation{}

	content := strings.Split(ctx.RawArgs, "\n")
	if len(content) < 2 {
		_, err = ctx.Sendf("Not enough arguments provided.")
		return err
	}
	names, err := shellwords.Parse(content[0])
	if err != nil {
		names = strings.Split(content[0], " ")
	}
	e.Name = names[0]
	if len(names) > 1 {
		e.Aliases = names[1:]
	}
	e.Description = strings.Join(content[1:], "\n")
	e, err = c.db.AddExplanation(e)
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}
	_, err = ctx.Sendf("Added explanation with ID %v.\nName: `%v`\nAliases: `%v`", e.ID, e.Name, strings.Join(e.Aliases, ", "))
	return err
}
