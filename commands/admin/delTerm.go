package admin

import (
	"strconv"

	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
)

func (c *commands) delTerm(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return ctx.CommandError(err)
	}

	t, err := c.db.GetTerm(id)
	if err != nil {
		return ctx.CommandError(err)
	}

	m, err := ctx.Send(&discordgo.MessageSend{
		Content: "Are you sure you want to delete this term? React with ✅ to delete it, or with ❌ to cancel.",
		Embed:   t.TermEmbed(""),
	})

	ctx.AddYesNoHandler(m.ID, func(ctx *crouter.Ctx) {
		err = c.db.RemoveTerm(id)
		if err != nil {
			ctx.CommandError(err)
			return
		}

		ctx.Sendf("Removed term `%v`.", id)
	}, func(ctx *crouter.Ctx) {
		ctx.Sendf("Cancelled.")
	})
	return
}
