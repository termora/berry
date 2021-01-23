package static

import (
	"math/rand"
	"time"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/pkgo"
)

var pk = pkgo.NewSession(nil)

var greetings = []string{"Hello", "Heyo", "Heya", "Hiya"}

func (c *Commands) hello(ctx *bcr.Context) (err error) {
	// sleep for a second to give PK time to process the message
	time.Sleep(1 * time.Second)

	var name string
	m, err := pk.GetMessage(ctx.Message.ID.String())
	if err != nil {
		member, err := ctx.ParseMember(ctx.Author.ID.String())
		if err != nil {
			name = ctx.Author.Username
		} else {
			if member.Nick != "" {
				name = member.Nick
			} else {
				name = ctx.Author.Username
			}
		}
	} else {
		name = m.Member.Name
	}

	_, err = ctx.Sendf("%v, %v!", greetings[rand.Intn(len(greetings))], name)
	return err
}
