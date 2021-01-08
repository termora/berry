package admin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Starshine113/bcr"
)

func (c *commands) update(ctx *bcr.Context) (err error) {
	_, err = ctx.Send("Updating Git repository...", nil)
	if err != nil {
		return err
	}

	git := exec.Command("git", "pull")
	pullOutput, err := git.Output()
	if err != nil {
		_, err = ctx.Send(fmt.Sprintf("Error pulling repository:\n```%v```", err), nil)
		return err
	}
	_, err = ctx.Send(fmt.Sprintf("Git:\n```%v```", string(pullOutput)), nil)
	if err != nil {
		return err
	}

	update := exec.Command("/usr/local/go/bin/go", "build")
	updateOutput, err := update.Output()
	if err != nil {
		_, err = ctx.Send(fmt.Sprintf("Error building:\n```%v```", err), nil)
		return err
	}
	_, err = ctx.Send(fmt.Sprintf("`go build`:\n```%v```", bcr.DefaultValue(string(updateOutput), "[no output]")), nil)
	return
}

func (c *commands) kill(ctx *bcr.Context) (err error) {
	_, err = ctx.Send("Restarting the bot, please wait...", nil)
	if err != nil {
		return err
	}
	c.sugar.Infof("Kill command received, shutting down...")

	ctx.Router.Session.Close()
	c.sugar.Infof("Disconnected from Discord.")
	c.db.Pool.Close()
	c.sugar.Infof("Closed database connection.")
	os.Exit(0)
	return nil
}
