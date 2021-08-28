package auditlog

import (
	"context"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

func (bot *AuditLog) reason(ctx *bcr.Context) (err error) {
	reason := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if reason == ctx.RawArgs {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	if reason == "" || len(reason) > 1000 {
		_, err = ctx.Replyc(bcr.ColourRed, "Empty reason, or it's too long.")
		return
	}

	var id int64
	if strings.EqualFold(ctx.Args[0], "l") || strings.EqualFold(ctx.Args[0], "latest") {
		err = bot.DB.QueryRow(context.Background(), "select coalesce(max(id), 0) from audit_log where user_id = $1", ctx.Author.ID).Scan(&id)
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		}
	} else {
		id, err = strconv.ParseInt(ctx.Args[0], 10, 0)
		if err != nil || id < 1 {
			return ctx.SendX("Couldn't parse your input as a number.")
		}
	}

	err = bot.DB.QueryRow(context.Background(), "select coalesce(id, 0) from audit_log where user_id = $1 and id = $2", ctx.Author.ID, id).Scan(&id)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return ctx.SendX("Either there is no audit log entry with that ID, or it's not yours.")
		}
	}

	entry, err := bot.updateReason(id, reason)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	if entry.PublicMessageID.IsValid() && bot.Config.Bot.AuditLog.Public.IsValid() {
		_, err = bot.State.EditEmbeds(bot.Config.Bot.AuditLog.Public, entry.PublicMessageID, bot.publicEmbed(entry, bot.desc(entry)))
		if err != nil {
			ctx.Replyc(bcr.ColourRed, "Couldn't edit public log entry.")
		}
	}
	if entry.PrivateMessageID.IsValid() && bot.Config.Bot.AuditLog.Private.IsValid() {
		_, err = bot.State.EditEmbeds(bot.Config.Bot.AuditLog.Private, entry.PrivateMessageID, bot.privateEmbeds(entry)...)
		if err != nil {
			ctx.Replyc(bcr.ColourRed, "Couldn't edit private log entry.")
		}
	}

	_, err = ctx.Replyc(bcr.ColourGreen, "Updated reason for entry %v!", entry.ID)
	return
}
