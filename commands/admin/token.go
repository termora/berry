package admin

import (
	"fmt"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) token(ctx *bcr.Context) (err error) {
	u, err := ctx.State.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		c.Sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?", nil)
		return
	}

	token, err := c.DB.GetOrCreateToken(ctx.Author.ID.String())
	if err != nil {
		c.Sugar.Errorf("Error creating token for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error creating/finding your token.", nil)
		return err
	}

	_, err = ctx.State.SendMessage(u.ID, fmt.Sprintf("⚠️ Please note that this token allows you to add, edit, and remove terms with the API. **Do not share this with anyone under any circumstances.**\nIf you lose your token, or it is expired, you can refresh it with `%vtoken refresh`. Your token is below:", ctx.Prefix), nil)
	if err != nil {
		return err
	}

	_, err = ctx.State.SendMessage(u.ID, token, nil)
	if err != nil {
		return err
	}

	if u.ID != ctx.Channel.ID {
		_, err = ctx.Send("✅ Check your DMs!", nil)
	}
	return err
}

func (c *Admin) refreshToken(ctx *bcr.Context) (err error) {
	u, err := ctx.State.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		c.Sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?", nil)
		return
	}

	token, err := c.DB.ResetToken(ctx.Author.ID.String())
	if err != nil {
		c.Sugar.Errorf("Error resetting token for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error resetting your token.", nil)
		return err
	}

	_, err = ctx.State.SendMessage(u.ID, fmt.Sprintf("⚠️ Please note that this token allows you to add, edit, and remove terms with the API. **Do not share this with anyone under any circumstances.**\nIf you lose your token, or it is expired, you can refresh it with `%vtoken refresh`. Your token is below:", ctx.Prefix), nil)
	if err != nil {
		return err
	}

	_, err = ctx.State.SendMessage(u.ID, token, nil)
	if err != nil {
		return err
	}

	if u.ID != ctx.Channel.ID {
		_, err = ctx.Send("✅ Check your DMs!", nil)
	}
	return err
}
