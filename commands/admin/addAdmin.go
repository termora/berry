package admin

import (
	"strings"

	"github.com/Starshine113/bcr"
)

func (c *commands) addAdmin(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to pass a user ID or mention to this command.", nil)
		return
	}

	u, err := ctx.ParseMember(strings.Join(ctx.Args, " "))
	if err != nil {
		_, err = ctx.Send("User not found", nil)
		return
	}

	msg, err := ctx.Sendf("Are you sure you want to add %v as a bot admin?", u.User.Mention())
	ctx.AddYesNoHandler(*msg, ctx.Author.ID, func(ctx *bcr.Context) {
		err := c.db.AddAdmin(u.User.ID.String())
		if err != nil {
			c.sugar.Errorf("Error adding admin %v: %v", u.User.ID.String(), err)
			_, err = ctx.Send("Error adding admin.", nil)
			if err != nil {
				c.sugar.Errorf("Error sending message: %v", err)
			}
			return
		}
		_, err = ctx.Sendf("Added %v as a bot admin.", u.Mention())
		if err != nil {
			c.sugar.Error("Error sending message:", err)
			return
		}
		c.admins, err = c.db.GetAdmins()
		if err != nil {
			c.sugar.Error("Error refreshing list of admins:", err)
		}
		return
	}, func(ctx *bcr.Context) {
		_, err = ctx.Send("Cancelled.", nil)
		if err != nil {
			c.sugar.Errorf("Error sending message: %v", err)
		}
		return
	})
	return err
}
