package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/utils/bot/extras/shellwords"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/admin/auditlog"
	"github.com/termora/berry/db"
)

func (c *Admin) addExplanation(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("Not enough arguments provided.")
		return err
	}
	e := &db.Explanation{}

	// split arguments--first line is aliases, second line on is explanation
	content := strings.Split(ctx.RawArgs, "\n")
	if len(content) < 2 {
		_, err = ctx.Sendf("Not enough arguments provided.")
		return err
	}

	// split the names
	names, err := shellwords.Parse(content[0])
	if err != nil {
		names = strings.Split(content[0], " ")
	}

	// name 0 is name, name 1+ is alias
	e.Name = names[0]
	if len(names) > 1 {
		e.Aliases = names[1:]
	}

	// and now we have to join the later strings again for the description
	e.Description = strings.Join(content[1:], "\n")

	// add the explanation
	e, err = c.DB.AddExplanation(e)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = c.AuditLog.SendLog(e.ID, auditlog.ExplanationEntry, auditlog.CreateAction, nil, e, ctx.Author.ID, nil)
	if err != nil {
		return
	}

	_, err = ctx.Sendf("Added explanation with ID %v.\nName: `%v`\nAliases: `%v`\n**:warning: Warning:** to have me respond to this explanation as a base command, please restart the bot.", e.ID, e.Name, strings.Join(e.Aliases, ", "))
	return err
}

func (c *Admin) toggleExplanationCmd(ctx *bcr.Context) (err error) {
	// we can't be bothered to check the *current* status, so just pass a value every time
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("Not enough arguments provided.")
		return err
	}

	// parse the ID
	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't parse a numeric ID.")
		return err
	}

	// parse the bool
	b, err := strconv.ParseBool(ctx.Args[1])
	if err != nil {
		_, err = ctx.Send("Couldn't parse a boolean.")
		return err
	}

	// set it in the database
	err = c.DB.SetAsCommand(id, b)
	if err != nil {
		_, err = ctx.Send(fmt.Sprintf("Internal error occurred: %v", bcr.AsCode(bcr.EscapeBackticks(err.Error()))))
		return err
	}

	_, err = ctx.Send(fmt.Sprintf("Set command status for `%v` to `%v`.\n**:warning: Note:** the bot has to be restarted to see the changes.", id, b))
	return err
}
