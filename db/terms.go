package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/georgysavva/scany/pgxscan"
)

// EmbedColour is the embed colour used throughout the bot
const EmbedColour = 0xc1302e

// Term holds info on a single term
type Term struct {
	ID           int
	Category     int
	CategoryName string
	Name         string
	Aliases      []string
	Description  string
	Source       string
	Created      time.Time

	// Rank is only populated with db.Search()
	Rank float64
	// Headline is only populated with db.Search()
	Headline string
}

// Errors related to database operations
var (
	ErrorNoRowsAffected = errors.New("no rows affected")
)

// Search searches the database for terms
func (db *Db) Search(input string, limit int) (terms []*Term, err error) {
	if limit == 0 {
		limit = 50
	}
	err = pgxscan.Select(context.Background(), db.Pool, &terms, `select
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.source, t.created,
	ts_rank_cd(t.searchtext, websearch_to_tsquery('english', $1), 32) as rank,
	ts_headline(t.description, websearch_to_tsquery('english', $1), 'StartSel=**, StopSel=**') as headline
	from public.terms as t, public.categories as c where t.searchtext @@ websearch_to_tsquery('english', $1) and t.category = c.id
	order by rank
	limit $2`, input, limit)
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
	t.id, t.category, c.name as category_name, t.name, t.aliases, t.description, t.source, t.created from public.terms as t, public.categories as c where t.id = $1`, id)
	return t, err
}

// TermEmbed creates a Discord embed from a term object
func (t *Term) TermEmbed() *discordgo.MessageEmbed {
	if t == nil {
		return nil
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	if len(t.Aliases) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Synonyms",
			Value: strings.Join(t.Aliases, ", "),
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "Source",
		Value: t.Source,
	})

	e := &discordgo.MessageEmbed{
		Title:       t.Name,
		Description: t.Description,
		Color:       EmbedColour,
		Timestamp:   t.Created.Format(time.RFC3339),
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("ID: %v | Category: %v (ID: %v) | Created", t.ID, t.CategoryName, t.Category),
		},
	}

	return e
}
