package static

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) ping(ctx bcr.Contexter) (err error) {
	t := time.Now()
	// this will return 0ms in the first minute after the bot is restarted
	// can't do much about that though
	heartbeat := ctx.Session().Gateway().EchoBeat().Sub(ctx.Session().Gateway().SentBeat())

	s := fmt.Sprintf("🏓 **Pong!**\nHeartbeat: %v", heartbeat)

	_, err = ctx.Send(s)
	if err != nil {
		return err
	}

	latency := time.Since(t).Round(time.Millisecond)

	_, err = ctx.EditOriginal(api.EditInteractionResponseData{
		Content: option.NewNullableString(
			fmt.Sprintf("%v\nMessage: %v", s, latency),
		),
	})
	return err
}
