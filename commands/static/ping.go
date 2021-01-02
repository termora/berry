package static

import (
	"time"

	"github.com/Starshine113/crouter"
)

func (c *Commands) ping(ctx *crouter.Ctx) (err error) {
	t := time.Now()
	ping := ctx.Session.HeartbeatLatency().Round(time.Millisecond)
	m, err := ctx.Sendf("🏓 **Pong!**\nHeartbeat: %s", ping)
	if err != nil {
		return err
	}
	latency := time.Since(t).Round(time.Millisecond)
	_, err = ctx.Editf(m, m.Content+"\nMessage: %s", latency)
	return err
}
