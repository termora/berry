package admin

import (
	"strconv"

	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/admin/auditlog"
)

func (bot *Bot) delTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckRequiredArgs(1); err != nil {
		_, err = ctx.Send("No term ID provided.")
		return err
	}

	id, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	t, err := bot.DB.GetTerm(id)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	m, err := ctx.Send("Are you sure you want to delete this term? React with ✅ to delete it, or with ❌ to cancel.", bot.DB.TermEmbed(t))
	if err != nil {
		return err
	}

	// confirm deleting the term
	if yes, timeout := ctx.YesNoHandler(*m, ctx.Author.ID); !yes || timeout {
		ctx.Send("Cancelled.")
		return
	}

	err = bot.DB.RemoveTerm(id)
	if err != nil {
		bot.Log.Error("Error removing term:", err)
		bot.DB.InternalError(ctx, err)
		return
	}

	_, err = bot.AuditLog.SendLog(t.ID, auditlog.TermEntry, auditlog.DeleteAction, t, nil, ctx.Author.ID, nil)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("✅ Term deleted.")
	if err != nil {
		bot.Log.Error("Error sending message:", err)
	}
	return nil
}
