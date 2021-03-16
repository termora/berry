package admin

import (
	"strings"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) addAdmin(ctx *bcr.Context) (err error) {
	// if there's no arguments, return
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to pass a user ID or mention to this command.", nil)
		return
	}

	// parse the member
	u, err := ctx.ParseUser(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.Send("User not found", nil)
		return
	}

	msg, err := ctx.Sendf("Are you sure you want to add %v as a bot admin?", u.Mention())

	// add a yes/no reaction handler
	if yes, timeout := ctx.YesNoHandler(*msg, ctx.Author.ID); !yes || timeout {
		_, err = ctx.Send("Cancelled.", nil)
		return
	}

	err = c.DB.AddAdmin(u.ID.String())
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Added %v as a bot admin.", u.Mention())

	// refresh the list of admins
	c.admins, err = c.DB.GetAdmins()
	if err != nil {
		c.Sugar.Error("Error refreshing list of admins:", err)
	}
	return

}
