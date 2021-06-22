package db

import (
	"context"
	"errors"
	"math/rand"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// EmbedColour is the embed colour used throughout the bot
const EmbedColour = 0xd14171

// Errors related to database operations
var (
	ErrorNoRowsAffected = errors.New("no rows affected")
)

// TermCount ...
func (db *Db) TermCount() (count int) {
	db.Pool.QueryRow(context.Background(), "select count(id) from public.terms").Scan(&count)
	return count
}

// GetTerms gets all terms not blocked by the given mask
func (db *Db) GetTerms(mask TermFlag) (terms []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = c.id
	order by t.name, t.id`, mask)
	return terms, err
}

// GetCategoryTerms gets terms by category
func (db *Db) GetCategoryTerms(id int, mask TermFlag) (terms []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = $2
	and t.category = c.id
	order by t.name, t.id`, mask, id)
	return terms, err
}

// Search searches the database for terms
func (db *Db) Search(input string, limit int, ignore []string) (terms []*Term, err error) {
	if limit == 0 {
		limit = 50
	}
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 8) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0
	and not $4 && tags
	order by rank desc
	limit $2`, input, limit, FlagSearchHidden, ignore)
	return terms, err
}

// TermName gets a term by name
func (db *Db) TermName(n string) (t []*Term, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c where (t.name ilike $1 or $2 ilike any(t.aliases)) and t.category = c.id`, n, n)
	return t, err
}

// SearchCat searches for terms from a single category
func (db *Db) SearchCat(input string, cat, limit int, showHidden bool, ignore []string) (terms []*Term, err error) {
	if limit == 0 {
		limit = 50
	}

	flags := FlagSearchHidden
	if showHidden {
		flags = 0
	}

	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 8) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c
	where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id and t.flags & $3 = 0 and t.category = $4
	and not $5 && tags
	order by rank desc
	limit $2`, input, limit, flags, cat, ignore)
	return terms, err
}

// AddTerm adds a term to the database
func (db *Db) AddTerm(t *Term) (*Term, error) {
	if t.Aliases == nil {
		t.Aliases = []string{}
	}
	if t.Tags == nil {
		t.Tags = []string{}
	}

	err := db.Pool.QueryRow(context.Background(), "insert into public.terms (name, category, aliases, description, source, aliases_string, tags) values ($1, $2, $3, $4, $5, $6, $7) returning id, created", t.Name, t.Category, t.Aliases, t.Description, t.Source, strings.Join(t.Aliases, ", "), t.Tags).Scan(&t.ID, &t.Created)
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
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c where t.id = $1 and t.category = c.id`, id)
	return t, err
}

// RandomTerm gets a random term from the database
func (db *Db) RandomTerm(ignore []string) (t *Term, err error) {
	var terms []*Term
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = c.id
	and not $2 && tags
	order by t.id`, FlagRandomHidden, ignore)
	if err != nil {
		return
	}

	if len(terms) == 1 {
		return terms[0], nil
	}

	if len(terms) == 0 {
		return nil, pgx.ErrNoRows
	}

	n := rand.Intn(len(terms) - 1)
	return terms[n], nil
}

// RandomTermCategory gets a random term from the database from the specified category
func (db *Db) RandomTermCategory(id int, ignore []string) (t *Term, err error) {
	var terms []*Term
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = c.id
	and t.category = $2
	and not $3 && tags
	order by t.id`, FlagRandomHidden, id, ignore)
	if err != nil {
		return
	}

	if len(terms) == 1 {
		return terms[0], nil
	}

	if len(terms) == 0 {
		return nil, pgx.ErrNoRows
	}

	n := rand.Intn(len(terms) - 1)
	return terms[n], nil
}

// SetFlags sets the flags for a term
func (db *Db) SetFlags(id int, flags TermFlag) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set flags = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", flags, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	return
}

// SetCW sets the content warning for a term
func (db *Db) SetCW(id int, text string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set content_warnings = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", text, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	return
}

// UpdateDesc updates the description for a term
func (db *Db) UpdateDesc(id int, desc string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set description = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", desc, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// UpdateSource updates the source for a term
func (db *Db) UpdateSource(id int, source string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set source = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", source, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// UpdateTitle updates the title for a term
func (db *Db) UpdateTitle(id int, title string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set name = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", title, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// UpdateImage updates the title for a term
func (db *Db) UpdateImage(id int, img string) (err error) {
	commandTag, err := db.Pool.Exec(context.Background(), "update public.terms set image_url = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", img, id)
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// UpdateAliases updates the aliases for a term
func (db *Db) UpdateAliases(id int, aliases []string) (err error) {
	var commandTag pgconn.CommandTag
	if len(aliases) > 0 {
		commandTag, err = db.Pool.Exec(context.Background(), "update public.terms set aliases = $1, aliases_string = $2, last_modified = (current_timestamp at time zone 'utc') where id = $3", aliases, strings.Join(aliases, ", "), id)
	} else {
		commandTag, err = db.Pool.Exec(context.Background(), "update public.terms set aliases = array[]::text[], aliases_string = '', last_modified = (current_timestamp at time zone 'utc') where id = $1", id)
	}
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// SetNote updates the note for a term
func (db *Db) SetNote(id int, note string) (err error) {
	var commandTag pgconn.CommandTag
	if len(note) > 0 {
		commandTag, err = db.Pool.Exec(context.Background(), "update public.terms set note = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", note, id)
	} else {
		commandTag, err = db.Pool.Exec(context.Background(), "update public.terms set note = '', last_modified = (current_timestamp at time zone 'utc') where id = $1", id)
	}
	if err != nil {
		return
	}
	if commandTag.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}
	return
}

// UpdateTags updates the tags for a term
func (db *Db) UpdateTags(id int, tags []string) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update public.terms set tags = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", tags, id)
	return
}
