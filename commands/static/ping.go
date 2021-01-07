package static

import (
	"fmt"
	"time"

	"github.com/Starshine113/bcr"
)

func (c *Commands) ping(ctx *bcr.Context) (err error) {
	t := time.Now()
	heartbeat := ctx.Session.Gateway.PacerLoop.EchoBeat.Time().Sub(ctx.Session.Gateway.PacerLoop.SentBeat.Time()).Round(time.Millisecond)

	m, err := ctx.Send(fmt.Sprintf("🏓 **Pong!**\nHeartbeat: %v", heartbeat), nil)
	if err != nil {
		return err
	}

	latency := time.Since(t).Round(time.Millisecond)

	_, err = ctx.Edit(m, fmt.Sprintf("%v\nMessage: %v", m.Content, latency), nil)
	return err
}
