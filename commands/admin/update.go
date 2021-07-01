package admin

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/starshine-sys/bcr"
)

func (c *Admin) update(ctx *bcr.Context) (err error) {
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

	output := string(pullOutput)
	if len(output) > 1900 {
		output = "...\n" + output[:len(output)-1900]
	}

	_, err = ctx.Send(fmt.Sprintf("Git:\n```%v```", output), nil)
	if err != nil {
		return err
	}

	t := time.Now()
	update := exec.Command("/usr/local/go/bin/go", "build", "-v")
	updateOutput, err := update.CombinedOutput()
	if err != nil {
		_, err = ctx.Send(fmt.Sprintf("Error building:\n```%v```", err), nil)
		return err
	}

	output = string(updateOutput)
	if len(output) > 1900 {
		output = "...\n" + output[:len(output)-1900]
	}
	if output == "" {
		output = "[no output]"
	}

	buildTime := time.Since(t).Round(time.Millisecond)
	_, err = ctx.Send(fmt.Sprintf("`go build` (%v):\n```%v```", buildTime, output), nil)
	return
}
