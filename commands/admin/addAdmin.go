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
	ctx.AddYesNoHandler(*msg, ctx.Author.ID, func(ctx *bcr.Context) {
		// add the admin
		err := c.DB.AddAdmin(u.ID.String())

		if err != nil {
			c.Sugar.Errorf("Error adding admin %v: %v", u.ID.String(), err)
			_, err = ctx.Send("Error adding admin.", nil)
			if err != nil {
				c.Sugar.Errorf("Error sending message: %v", err)
			}
			return
		}

		_, err = ctx.Sendf("Added %v as a bot admin.", u.Mention())
		if err != nil {
			c.Sugar.Error("Error sending message:", err)
			return
		}

		// refresh the list of admins
		c.admins, err = c.DB.GetAdmins()
		if err != nil {
			c.Sugar.Error("Error refreshing list of admins:", err)
		}
		return
	}, func(ctx *bcr.Context) {
		// otherwise cancel
		_, err = ctx.Send("Cancelled.", nil)
		if err != nil {
			c.Sugar.Errorf("Error sending message: %v", err)
		}
		return
	})
	return err
}
