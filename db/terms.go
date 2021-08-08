package db

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/termora/berry/db/search"
)

// EmbedColour is the embed colour used throughout the bot
const EmbedColour = 0xd14171

// Errors related to database operations
var (
	ErrorNoRowsAffected = errors.New("no rows affected")
)

// TermCount ...
func (db *Db) TermCount() (count int) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting term count")

	db.Pool.QueryRow(ctx, "select count(id) from public.terms").Scan(&count)
	return count
}

// GetTerms gets all terms not blocked by the given mask
func (db *Db) GetTerms(mask search.TermFlag) (terms []*Term, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting terms matching flags %v", mask)

	err = pgxscan.Select(ctx, db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = c.id
	order by t.name, t.id`, mask)
	return terms, err
}

// GetCategoryTerms gets terms by category
func (db *Db) GetCategoryTerms(id int, mask search.TermFlag) (terms []*Term, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting terms in category %v matching flags %v", id, mask)

	err = pgxscan.Select(ctx, db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.flags, t.tags, t.content_warnings, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = $2
	and t.category = c.id
	order by t.name, t.id`, mask, id)
	return terms, err
}

// TermName gets a term by name
func (db *Db) TermName(n string) (t []*Term, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting term with name %v", n)

	err = pgxscan.Select(ctx, db.Pool, &t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c where (t.name ilike $1 or $2 ilike any(t.aliases)) and t.category = c.id`, n, n)
	return t, err
}

// AddTerm adds a term to the database
func (db *Db) AddTerm(t *Term) (*Term, error) {
	if t.Aliases == nil {
		t.Aliases = []string{}
	}
	if t.Tags == nil {
		t.Tags = []string{}
	}

	Debug("Adding term %v", t.Name)

	ctx, cancel := db.Context()
	defer cancel()

	err := db.Pool.QueryRow(ctx, "insert into public.terms (name, category, aliases, description, source, aliases_string, tags) values ($1, $2, $3, $4, $5, $6, $7) returning id, created", t.Name, t.Category, t.Aliases, t.Description, t.Source, strings.Join(t.Aliases, ", "), t.Tags).Scan(&t.ID, &t.Created)
	if err != nil {
		return nil, err
	}

	return t, db.SyncTerm(t)
}

// RemoveTerm removes a term from the database
func (db *Db) RemoveTerm(id int) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Deleting term %v", id)

	ct, err := db.Pool.Exec(ctx, "delete from public.terms where id = $1", id)
	if err != nil {
		return
	}
	if ct.RowsAffected() != 1 {
		return ErrorNoRowsAffected
	}

	return db.SyncDelete(id)
}

// GetTerm gets a term by ID
func (db *Db) GetTerm(id int) (t *Term, err error) {
	t = &Term{}

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting term %v", id)

	err = pgxscan.Get(ctx, db.Pool, t, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags, t.image_url,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c where t.id = $1 and t.category = c.id`, id)
	return t, err
}

// RandomTerm gets a random term from the database
func (db *Db) RandomTerm(ignore []string) (t *Term, err error) {
	var terms []*Term

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting random term ignoring `%v`", ignore)

	err = pgxscan.Select(ctx, db.Pool, &terms, `select t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = c.id
	and not $2 && tags
	order by t.id`, search.FlagRandomHidden, ignore)
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

	ctx, cancel := db.Context()
	defer cancel()

	Debug("Getting random term in %v ignoring `%v`", id, ignore)

	err = pgxscan.Select(ctx, db.Pool, &terms, `select t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.note, t.source, t.created, t.last_modified, t.content_warnings, t.flags, t.tags,
	array(select display from public.tags where normalized = any(t.tags)) as display_tags
	from public.terms as t, public.categories as c
	where t.flags & $1 = 0 and t.category = c.id
	and t.category = $2
	and not $3 && tags
	order by t.id`, search.FlagRandomHidden, id, ignore)
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
func (db *Db) SetFlags(id int, flags search.TermFlag) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Setting flags for %v to %v", id, flags)

	commandTag, err := db.Pool.Exec(ctx, "update public.terms set flags = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", flags, id)
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
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Setting cw for %v to `%v`", id, text)

	commandTag, err := db.Pool.Exec(ctx, "update public.terms set content_warnings = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2", text, id)
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
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating description for %v to `%v`", id, desc)

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set description = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2 returning id, name, category, aliases, description, source, tags", desc, id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}

// UpdateSource updates the source for a term
func (db *Db) UpdateSource(id int, source string) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating source for %v to `%v`", id, source)

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set source = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2 returning id, name, category, aliases, description, source, tags", source, id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}

// UpdateTitle updates the title for a term
func (db *Db) UpdateTitle(id int, title string) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating title for %v to `%v`", id, title)

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set name = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2 returning id, name, category, aliases, description, source, tags", title, id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}

// UpdateImage updates the image for a term
func (db *Db) UpdateImage(id int, img string) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating image for %v to `%v`", id, img)

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set image_url = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2 returning id, name, category, aliases, description, source, tags", img, id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}

// UpdateAliases updates the aliases for a term
func (db *Db) UpdateAliases(id int, aliases []string) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating aliases for %v to `%v`", id, aliases)

	if aliases == nil {
		aliases = []string{}
	}

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set aliases = $1, aliases_string = $2, last_modified = (current_timestamp at time zone 'utc') where id = $3 returning id, name, category, aliases, description, source, tags", aliases, strings.Join(aliases, ", "), id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}

// SetNote updates the note for a term
func (db *Db) SetNote(id int, note string) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating note for %v to `%v`", id, note)

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set note = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2 returning id, name, category, aliases, description, source, tags", note, id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}

// UpdateTags updates the tags for a term
func (db *Db) UpdateTags(id int, tags []string) (err error) {
	ctx, cancel := db.Context()
	defer cancel()

	Debug("Updating tags for %v to `%v`", id, tags)

	var t Term
	err = pgxscan.Get(ctx, db.Pool, &t, "update public.terms set tags = $1, last_modified = (current_timestamp at time zone 'utc') where id = $2 returning id, name, category, aliases, description, source, tags", tags, id)
	if err != nil {
		return
	}

	return db.SyncTerm(&t)
}
