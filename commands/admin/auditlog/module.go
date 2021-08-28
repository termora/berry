package auditlog

import (
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/bot"
)

// AuditLog ...
type AuditLog struct {
	State *state.State

	*bot.Bot
}

// New returns a new AuditLog object.
func New(bot *bot.Bot) *AuditLog {
	st, _ := bot.Router.StateFromGuildID(0)

	return &AuditLog{
		State: st,
		Bot:   bot,
	}
}

// Init ...
func Init(bot *bot.Bot, perms bcr.CustomPerms) (s string, list []*bcr.Command) {
	s = "Audit log"

	b := New(bot)

	_ = b

	b.Router.AddCommand(&bcr.Command{
		Name:              "reason",
		Summary:           "Set a reason for the given audit log entry.",
		Usage:             "<ID|latest> <reason...>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: perms,
		Command:           b.reason,
		Hidden:            true,
	})

	cmd := b.Router.GetCommand("admin")

	cmd.AddSubcommand(&bcr.Command{
		Name:              "reason",
		Summary:           "Set a reason for the given audit log entry.",
		Usage:             "<ID|latest> <reason...>",
		Args:              bcr.MinArgs(2),
		CustomPermissions: perms,
		Command:           b.reason,
	})

	return
}
