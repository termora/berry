package db

import (
	"context"
	"errors"
	"math/rand"

	"github.com/georgysavva/scany/pgxscan"
)

// EmbedColour is the embed colour used throughout the bot
const EmbedColour = 0xe00d7a

// Errors related to database operations
var (
	ErrorNoRowsAffected = errors.New("no rows affected")
)

// GetTerms gets all terms not blocked by the given mask
func (db *Db) GetTerms(mask TermFlag) (terms []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.source, t.created, t.flags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0
	order by t.id`, mask)
	return terms, err
}

// Search searches the database for terms
func (db *Db) Search(input string, limit int) (terms []*Term, err error) {
	if limit == 0 {
		limit = 50
	}
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.source, t.created, t.flags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 32) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0
	order by rank desc
	limit $2`, input, limit, FlagSearchHidden)
	return terms, err
}

// AddTerm adds a term to the database
func (db *Db) AddTerm(t *Term) (*Term, error) {
	err := db.Pool.QueryRow(context.Background(), "insert into public.terms (name, category, aliases, description, source) values ($1, $2, $3, $4, $5) returning id, created", t.Name, t.Category, t.Aliases, t.Description, t.Source).Scan(&t.ID, &t.Created)
	return t, err
}

// RemoveTerm removes a term from the database
func (db *Db) RemoveTerm(id int) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "delete from public.terms where id = $1", id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// GetTerm gets a term by ID
func (db *Db) GetTerm(id int) (t *Term, err error) {
	t = &Term{}
	err = pgxscan.Get(context.Background(), db.Pool, t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.source, t.created, t.flags from public.terms as t, public.categories as c where t.id = $1`, id)
	return t, err
}

// RandomTerm gets a random term from the database
func (db *Db) RandomTerm() (t *Term, err error) {
	var terms []*Term
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.source, t.created, t.flags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0
	order by t.id`, FlagRandomHidden)
	if err != nil {
		return
	}

	if len(terms) == 1 {
		return terms[0], nil
	}

	n := rand.Intn(len(terms) - 1)
	return terms[n], nil
}

// SetFlags sets the flags for a term
func (db *Db) SetFlags(id int, flags TermFlag) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set flags = $1 where id = $2", flags, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	return
}
