package db

import (
	"context"
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
func (db *Db) AddExplanation(e *Explanation) (ex *Explanation, err error) {
	err = db.Pool.QueryRow(context.Background(), "insert into public.explanations (name, aliases, description) values ($1, $2, $3) returning id, created", e.Name, e.Aliases, e.Description).Scan(&e.ID, &e.Created)
	return e, err
}

// GetExplanation ...
func (db *Db) GetExplanation(s string) (e *Explanation, err error) {
	e = &Explanation{}
	err = pgxscan.Get(context.Background(), db.Pool, e, "select id, name, aliases, description, created, as_command from public.explanations where lower(name) = lower($1) order by id desc limit 1", s)
	return e, err
}

// GetAllExplanations ...
func (db *Db) GetAllExplanations() (e []*Explanation, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &e, "select id, name, aliases, description, created, as_command from public.explanations order by id")
	return e, err
}

// GetCmdExplanations ...
func (db *Db) GetCmdExplanations() (e []*Explanation, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &e, "select id, name, aliases, description, created, as_command from public.explanations where as_command = true order by id")
	return e, err
}

// SetAsCommand ...
func (db *Db) SetAsCommand(id int, b bool) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.explanations set as_command = $1 where id = $2", b, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}
