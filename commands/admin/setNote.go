package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Admin) setNote(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		var notes string
		for k := range c.Config.QuickNotes {
			notes += fmt.Sprintf("`%v`\n", k)
		}
		if notes == "" {
			notes = "No quick notes."
		}

		_, err = ctx.Send("Not enough arguments provided: need ID and note (or \"clear\")", &discord.Embed{
			Title:       "List of quick notes",
			Description: notes,
			Color:       db.EmbedColour,
		})
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	t, err := c.DB.GetTerm(id)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	note := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	// if the input is "clear", remove the note
	if note == "clear" || note == "-clear" {
		note = ""
	}

	if n, ok := c.Bot.Config.QuickNotes[note]; ok && n != "" {
		note = n
	}

	if len(note) > 1000 {
		_, err = ctx.Sendf("❌ The note you gave is too long (%v > 1000 characters).", len(note))
		return
	}

	err = c.DB.SetNote(t.ID, note)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send(fmt.Sprintf("Updated note for %v.", id), &discord.Embed{
		Description: note,
		Color:       db.EmbedColour,
	})
	return
}
