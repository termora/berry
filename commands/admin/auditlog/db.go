package auditlog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/termora/berry/db"
)

// EntrySubject ...
type EntrySubject string

// ...
const (
	TermEntry        EntrySubject = "term"
	PronounsEntry    EntrySubject = "pronouns"
	ExplanationEntry EntrySubject = "explanation"
)

// ActionType ...
type ActionType string

// ...
const (
	CreateAction ActionType = "create"
	UpdateAction ActionType = "update"
	DeleteAction ActionType = "delete"
)

// ...
const (
	ErrBeforeNil = errors.Sentinel("before is nil")
	ErrAfterNil  = errors.Sentinel("after is nil")

	ErrInvalidSubjectType = errors.Sentinel("invalid subject")
)

// Entry ...
type Entry struct {
	ID int64

	SubjectID int
	Subject   EntrySubject
	Action    ActionType

	Before json.RawMessage
	After  json.RawMessage

	UserID discord.UserID
	Reason sql.NullString

	Timestamp time.Time

	PublicMessageID  discord.MessageID
	PrivateMessageID discord.MessageID
}

// BeforeTerm gets the entry's Before as a db.Term
func (e *Entry) BeforeTerm() (t db.Term, err error) {
	if e.Subject != TermEntry {
		return t, ErrInvalidSubjectType
	}

	if e.Action == CreateAction {
		return e.AfterTerm()
	}

	if e.Before == nil {
		return t, ErrBeforeNil
	}

	err = json.Unmarshal(e.Before, &t)
	return
}

// BeforePronouns gets the entry's Before as a db.PronounSet
func (e *Entry) BeforePronouns() (p db.PronounSet, err error) {
	if e.Subject != PronounsEntry {
		return p, ErrInvalidSubjectType
	}

	if e.Action == CreateAction {
		return e.AfterPronouns()
	}

	if e.Before == nil {
		return p, ErrBeforeNil
	}

	err = json.Unmarshal(e.Before, &p)
	return
}

// BeforeExplanation gets the entry's Before as a db.Explanation
func (e *Entry) BeforeExplanation() (ex db.Explanation, err error) {
	if e.Subject != ExplanationEntry {
		return ex, ErrInvalidSubjectType
	}

	if e.Action == CreateAction {
		return e.AfterExplanation()
	}

	if e.Before == nil {
		return ex, ErrBeforeNil
	}

	err = json.Unmarshal(e.Before, &ex)
	return
}

// AfterTerm gets the entry's After as a db.Term
func (e *Entry) AfterTerm() (t db.Term, err error) {
	if e.Subject != TermEntry {
		return t, ErrInvalidSubjectType
	}

	if e.Action == DeleteAction {
		return e.BeforeTerm()
	}

	if e.After == nil {
		return t, ErrAfterNil
	}

	err = json.Unmarshal(e.After, &t)
	return
}

// AfterPronouns gets the entry's After as a db.PronounSet
func (e *Entry) AfterPronouns() (p db.PronounSet, err error) {
	if e.Subject != PronounsEntry {
		return p, ErrInvalidSubjectType
	}

	if e.Action == DeleteAction {
		return e.BeforePronouns()
	}

	if e.After == nil {
		return p, ErrAfterNil
	}

	err = json.Unmarshal(e.After, &p)
	return
}

// AfterExplanation gets the entry's After as a db.Explanation
func (e *Entry) AfterExplanation() (ex db.Explanation, err error) {
	if e.Subject != ExplanationEntry {
		return ex, ErrInvalidSubjectType
	}

	if e.Action == DeleteAction {
		return e.BeforeExplanation()
	}

	if e.After == nil {
		return ex, ErrAfterNil
	}

	err = json.Unmarshal(e.After, &ex)
	return
}

func (bot *AuditLog) insertEntry(subjectID int, subjectType EntrySubject, actionType ActionType, before, after interface{}, userID discord.UserID, reason *string) (e Entry, err error) {
	s := sql.NullString{Valid: false}
	if reason != nil {
		s = sql.NullString{
			Valid:  true,
			String: *reason,
		}
	}

	err = pgxscan.Get(context.Background(), bot.DB.Pool, &e, "insert into audit_log (subject_id, subject, action, before, after, user_id, reason) values ($1, $2, $3, $4, $5, $6, $7) returning *", subjectID, subjectType, actionType, before, after, userID, s)
	return
}

// ErrAuditLogSendFailed ...
const ErrAuditLogSendFailed = errors.Sentinel("failed to send one or both audit log entries")

func (bot *AuditLog) desc(entry Entry) string {
	desc := strings.Title(string(entry.Action)) + "d"

	switch entry.Subject {
	case TermEntry:
		term, err := entry.BeforeTerm()
		if err != nil {
			bot.DB.Sugar.Errorf("Error unmarshaling term: %v", err)
		}

		desc += " term **"
		if entry.Action != DeleteAction && bot.DB.TermBaseURL != "" {
			desc += fmt.Sprintf("[%v](%v%v)", term.Name, bot.DB.TermBaseURL, term.ID)
		} else {
			desc += term.Name
		}
		desc += "**"
	case PronounsEntry:
		p, _ := entry.BeforePronouns()

		desc += " pronouns **" + p.String() + "**"
	case ExplanationEntry:
		ex, _ := entry.BeforeExplanation()

		desc += " explanation **" + ex.Name + "**"
	}
	return desc
}

// SendLog sends an audit log
func (bot *AuditLog) SendLog(subjectID int, subjectType EntrySubject, actionType ActionType, before, after interface{}, userID discord.UserID, reason *string) (id int64, err error) {
	entry, err := bot.insertEntry(subjectID, subjectType, actionType, before, after, userID, reason)
	if err != nil {
		return 0, err
	}

	desc := bot.desc(entry)
	publicID, err := bot.sendPublicEmbed(entry, desc)
	if err != nil {
		return entry.ID, errors.Wrap(err, string(ErrAuditLogSendFailed))
	}
	privateID, err := bot.sendPrivateEmbed(entry)
	if err != nil {
		return entry.ID, errors.Wrap(err, string(ErrAuditLogSendFailed))
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update audit_log set public_message_id = $1, private_message_id = $2 where id = $3", publicID, privateID, entry.ID)
	return entry.ID, err
}

func (bot *AuditLog) updateReason(id int64, reason string) (e Entry, err error) {
	s := sql.NullString{
		Valid:  true,
		String: reason,
	}

	err = pgxscan.Get(context.Background(), bot.DB.Pool, &e, "update audit_log set reason = $1 where id = $2 returning *", s, id)
	return
}
