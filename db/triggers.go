package db

import (
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

// Explanation is a single explanation
type Explanation struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Aliases     []string  `json:"aliases"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`

	AsCommand bool `json:"-"`
}

// AddExplanation adds an explanation to the database
func (db *DB) AddExplanation(e *Explanation) (ex *Explanation, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Adding explanation %v with description %v", e.Name, e.Description)

	err = db.QueryRow(ctx, "insert into public.explanations (name, aliases, description) values ($1, $2, $3) returning id, created", e.Name, e.Aliases, e.Description).Scan(&e.ID, &e.Created)
	return e, err
}

// GetExplanation ...
func (db *DB) GetExplanation(s string) (e *Explanation, err error) {
	e = &Explanation{}

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting explanation with name %v", s)

	err = pgxscan.Get(ctx, db.Pool, e, "select id, name, aliases, description, created, as_command from public.explanations where lower(name) = lower($1) order by id desc limit 1", s)
	return e, err
}

// GetAllExplanations ...
func (db *DB) GetAllExplanations() (e []*Explanation, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting all explanations")

	err = pgxscan.Select(ctx, db.Pool, &e, "select id, name, aliases, description, created, as_command from public.explanations order by id")
	return e, err
}

// GetCmdExplanations ...
func (db *DB) GetCmdExplanations() (e []*Explanation, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting command explanations")

	err = pgxscan.Select(ctx, db.Pool, &e, "select id, name, aliases, description, created, as_command from public.explanations where as_command = true order by id")
	return e, err
}

// SetAsCommand ...
func (db *DB) SetAsCommand(id int, b bool) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Setting %v as command: %v", id, b)

	commandTag, err := db.Exec(ctx, "update public.explanations set as_command = $1 where id = $2", b, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}
