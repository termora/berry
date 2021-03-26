package static

import (
	"fmt"
	"time"

	"github.com/starshine-sys/bcr"
)

func (c *Commands) ping(ctx *bcr.Context) (err error) {
	t := time.Now()
	// this will return 0ms in the first minute after the bot is restarted
	// can't do much about that though
	heartbeat := ctx.State.Gateway.PacerLoop.EchoBeat.Time().Sub(ctx.State.Gateway.PacerLoop.SentBeat.Time()).Round(time.Millisecond)

	m, err := ctx.Send(fmt.Sprintf("🏓 **Pong!**\nHeartbeat: %v", heartbeat), nil)
	if err != nil {
		return err
	}

	latency := time.Since(t).Round(time.Millisecond)

	_, err = ctx.Edit(m, fmt.Sprintf("%v\nMessage: %v", m.Content, latency), nil)
	return err
}
