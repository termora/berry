package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
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
		return c.db.InternalError(ctx, err)
	}
	_, err = ctx.Sendf("Added explanation with ID %v.\nName: `%v`\nAliases: `%v`\n**:warning: Warning:** to have me respond to this explanation as a base command, please restart the bot.", e.ID, e.Name, strings.Join(e.Aliases, ", "))
	return err
}

func (c *commands) toggleExplanationCmd(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided.", nil)
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't parse a numeric ID.", nil)
		return err
	}

	b, err := strconv.ParseBool(ctx.Args[1])
	if err != nil {
		_, err = ctx.Send("Couldn't parse a boolean.", nil)
		return err
	}

	err = c.db.SetAsCommand(id, b)
	if err != nil {
		_, err = ctx.Send(fmt.Sprintf("Internal error occurred: %v", bcr.AsCode(bcr.EscapeBackticks(err.Error()))), nil)
		return err
	}

	_, err = ctx.Send(fmt.Sprintf("Set command status for `%v` to `%v`.\n**:warning: Note:** the bot has to be restarted to see the changes.", id, b), nil)
	return err
}
