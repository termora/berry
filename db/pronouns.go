package db

import (
	"errors"
	"math/rand"
	"strings"

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

func (p PronounSet) String() string {
	return p.Subjective + "/" + p.Objective + "/" + p.PossDet + "/" + p.PossPro + "/" + p.Reflexive
}

// Errors ...
var (
	ErrMoreThanOneRow = errors.New("more than one row returned")
	ErrTooManyForms   = errors.New("too many forms given")
	ErrNoForms        = errors.New("no forms given")
)

// GetPronoun gets a pronoun from the database
// gods this function is shit but idc, if it works it works
func (db *Db) GetPronoun(forms ...string) (sets []*PronounSet, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting pronouns %v", strings.Join(forms, "/"))

	switch len(forms) {
	case 0:
		return nil, ErrNoForms
	case 1:
		err = pgxscan.Select(ctx, db.Pool, &sets, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where lower(subjective) = lower($1) order by sorting, subjective, objective, poss_det, poss_pro, reflexive", forms[0])
		if err != nil {
			return
		}
	case 2:
		err = pgxscan.Select(ctx, db.Pool, &sets, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where lower(subjective) = lower($1) and lower(objective) = lower($2) order by sorting, subjective, objective, poss_det, poss_pro, reflexive", forms[0], forms[1])
		if err != nil {
			return
		}
	case 3:
		err = pgxscan.Select(ctx, db.Pool, &sets, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where lower(subjective) = lower($1) and lower(objective) = lower($2) and lower(poss_det) = lower($3) order by sorting, subjective, objective, poss_det, poss_pro, reflexive", forms[0], forms[1], forms[2])
		if err != nil {
			return
		}
	case 4:
		err = pgxscan.Select(ctx, db.Pool, &sets, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where lower(subjective) = lower($1) and lower(objective) = lower($2) and lower(poss_det) = lower($3) and lower(poss_pro) = lower($4) order by sorting, subjective, objective, poss_det, poss_pro, reflexive", forms[0], forms[1], forms[2], forms[3])
		if err != nil {
			return
		}
	case 5:
		err = pgxscan.Select(ctx, db.Pool, &sets, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns where lower(subjective) = lower($1) and lower(objective) = lower($2) and lower(poss_det) = lower($3) and lower(poss_pro) = lower($4) and lower(reflexive) = lower($5) order by sorting, subjective, objective, poss_det, poss_pro, reflexive", forms[0], forms[1], forms[2], forms[3], forms[4])
		if err != nil {
			return
		}
	default:
		return nil, ErrTooManyForms
	}
	if len(sets) == 0 {
		return nil, pgx.ErrNoRows
	}
	return sets, nil
}

// RandomPronouns gets a random pronoun set from the database
func (db *Db) RandomPronouns() (p *PronounSet, err error) {
	var pronouns []*PronounSet

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting random pronouns")

	err = pgxscan.Select(ctx, db.Pool, &pronouns, `select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns order by id`)
	if err != nil {
		return
	}

	if len(pronouns) == 1 {
		return pronouns[0], nil
	}

	n := rand.Intn(len(pronouns) - 1)
	return pronouns[n], nil
}

// AddPronoun adds a pronoun set, returning the ID
func (db *Db) AddPronoun(p PronounSet) (id int, err error) {
	if p.Subjective == "" || p.Objective == "" || p.PossDet == "" || p.PossPro == "" || p.Reflexive == "" {
		return 0, ErrNoForms
	}

	Debug("Adding pronouns %s", p)

	ctx, cancel := db.Context()
	defer cancel()

	err = db.QueryRow(ctx, "insert into pronouns (subjective, objective, poss_det, poss_pro, reflexive) values ($1, $2, $3, $4, $5) returning id",
		strings.TrimSpace(p.Subjective), strings.TrimSpace(p.Objective), strings.TrimSpace(p.PossDet), strings.TrimSpace(p.PossPro), strings.TrimSpace(p.Reflexive),
	).Scan(&id)
	return id, err
}

// Pronouns ...
func (db *Db) Pronouns() (p []*PronounSet, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting all pronouns")

	err = pgxscan.Select(ctx, db.Pool, &p, "select id, subjective, objective, poss_det, poss_pro, reflexive from pronouns order by sorting, subjective, objective, poss_det, poss_pro, reflexive")
	return
}
