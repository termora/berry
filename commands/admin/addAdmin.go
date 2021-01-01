package admin

import (
	"strings"

	"github.com/Starshine113/crouter"
)

func (c *commands) addAdmin(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to pass a user ID or mention to this command.")
		return
	}

	u, err := ctx.ParseUser(strings.Join(ctx.Args, " "))
	if err != nil {
		if err == crouter.ErrNoID {
			_, err = ctx.Send("Invalid ID passed.")
		} else {
			_, err = ctx.Send("User not found")
		}
		return
	}

	msg, err := ctx.Sendf("Are you sure you want to add %v as a bot admin?", u.Mention())
	ctx.AddYesNoHandler(msg.ID, func(ctx *crouter.Ctx) {
		err := c.db.AddAdmin(u.ID)
		if err != nil {
			c.sugar.Errorf("Error adding admin %v: %v", u.ID, err)
			_, err = ctx.Send("Error adding admin.")
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
	}, func(ctx *crouter.Ctx) {
		_, err = ctx.Send("Cancelled.")
		if err != nil {
			c.sugar.Errorf("Error sending message: %v", err)
		}
		return
	})
	return err
}
