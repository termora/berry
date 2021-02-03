package db

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

// PronounSet is a single set of pronouns
type PronounSet struct {
	ID         int    `json:"id"`
	Subjective string `json:"subjective"`
	Objective  string `json:"objective"`
	PossDet    string `json:"possessive_determiner"`
	PossPro    string `json:"possessive_pronoun"`
	Reflexive  string `json:"reflexive"`
}

// Errors ...
var (
	ErrMoreThanOneRow = errors.New("more than one row returned")
	ErrTooManyForms   = errors.New("too many forms given")
	ErrNoForms        = errors.New("no forms given")
)

// GetPronoun gets a pronoun from the database
// gods this function is shit but idc, if it works it works
func (db *Db) GetPronoun(forms ...string) (set *PronounSet, err error) {
	var p []*PronounSet

	switch len(forms) {
	case 0:
		return nil, ErrNoForms
	case 1:
		err = pgxscan.Select(context.Background(), db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where subjective = $1 order by id", forms[0])
		if err != nil {
			return
		}
	case 2:
		err = pgxscan.Select(context.Background(), db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where subjective = $1 and objective = $2 order by id", forms[0], forms[1])
		if err != nil {
			return
		}
	case 3:
		err = pgxscan.Select(context.Background(), db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where subjective = $1 and objective = $2 and poss_det = $3 order by id", forms[0], forms[1], forms[2])
		if err != nil {
			return
		}
	case 4:
		err = pgxscan.Select(context.Background(), db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where subjective = $1 and objective = $2 and poss_det = $3 and poss_pro = $4 order by id", forms[0], forms[1], forms[2], forms[3])
		if err != nil {
			return
		}
	case 5:
		err = pgxscan.Select(context.Background(), db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where subjective = $1 and objective = $2 and poss_det = $3 and poss_pro = $4 and reflexive = $5 order by id", forms[0], forms[1], forms[2], forms[3], forms[4])
		if err != nil {
			return
		}
	default:
		return nil, ErrTooManyForms
	}

	if len(p) == 0 {
		return nil, pgx.ErrNoRows
	}

	if len(p) > 1 {
		return nil, ErrMoreThanOneRow
	}
	return p[0], nil
}

// AddPronoun adds a pronoun set, returning the ID
func (db *Db) AddPronoun(p PronounSet) (id int, err error) {
	if p.Subjective == "" || p.Objective == "" || p.PossDet == "" || p.PossPro == "" || p.Reflexive == "" {
		return 0, ErrNoForms
	}

	err = db.Pool.QueryRow(context.Background(), "insert into pronouns (subjective, objective, poss_det, poss_pro, reflexive) values ($1, $2, $3, $4, $5) returning id", p.Subjective, p.Objective, p.PossDet, p.PossPro, p.Reflexive).Scan(&id)
	return id, err
}

// Pronouns ...
func (db *Db) Pronouns() (p []*PronounSet, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns order by subjective")
	return
}
