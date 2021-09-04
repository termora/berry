package db

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/termora/berry/db/search"
)

// Term is an alias to search.Term
type Term = search.Term

// TermEmbed creates a Discord embed from a term object
func (db *Db) TermEmbed(t *Term) discord.Embed {
	if t == nil {
		return discord.Embed{Color: EmbedColour}
	}

	Debug("Creating term embed for %v", t.ID)

	e := discord.Embed{
		Title:     t.Name,
		Color:     EmbedColour,
		Timestamp: discord.NewTimestamp(t.Created),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Category: %v (ID: %v) | Created", t.ID, t.CategoryName, t.Category),
		},
	}

	var (
		desc = t.Description
		cw   = t.ContentWarnings
		note = t.Note
	)

	if db.TermBaseURL != "" {
		desc = db.LinkTerms(desc)
		note = db.LinkTerms(note)
		cw = db.LinkTerms(cw)
	}

	if cw != "" {
		desc = "||" + desc + "||"

		if len(desc) < 1024 {
			e.Description = fmt.Sprintf("**Content warning: %v**", cw)
		} else {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "​",
				Value: fmt.Sprintf("**Content warning: %v**", cw),
			})
		}
	}

	if len(desc) < 1024 && cw != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Description",
			Value: desc,
		})
	} else {
		e.Description = desc
	}

	if len(t.Aliases) != 0 {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Synonyms",
			Value: strings.Join(t.Aliases, ", "),
		})
	}

	if note != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Note",
			Value: note,
		})
	}

	if t.Warning() {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Warning",
			Value: "This term is only in this glossary for the sake of completeness. It may be derogatory, exclusionary, or harmful, especially when applied to other people and not as a self-description. Use this term with extreme caution.",
		})
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name:  "Source",
		Value: t.Source,
	})

	if len(t.DisplayTags) > 0 {
		var b strings.Builder
		for i, tag := range t.DisplayTags {
			if b.Len() >= 500 {
				b.WriteString(fmt.Sprintf("\nToo many to list (showing %v/%v)", i, len(t.DisplayTags)))
				break
			}
			b.WriteString(tag)
			if i != len(t.DisplayTags)-1 {
				b.WriteString(", ")
			}
		}

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Tag(s)",
			Value: b.String(),
		})
	}

	if db.TermBaseURL != "" {
		e.URL = db.TermBaseURL + strconv.Itoa(t.ID)
	}

	if t.ImageURL != "" {
		e.Image = &discord.EmbedImage{
			URL: t.ImageURL,
		}
	}

	return e
}

var linkRegexp = regexp.MustCompile(`\[\[(.*?)(\|.*?)?\]\]`)
var lowercaseRegexp = regexp.MustCompile(`[a-z]`)

// LinkTerms creates a strings.Replacer for all links in the page
func (db *Db) LinkTerms(input string) string {
	ctx, cancel := db.Context()
	defer cancel()

	// grab a single connection to use for the entire loop below
	// might not be any more performant than what we do normally but it doesn't hurt either so
	// ¯\_(ツ)_/¯
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return input
	}
	defer conn.Release()

	s := []string{}
	matches := linkRegexp.FindAllStringSubmatch(input, -1)

	for _, i := range matches {
		if len(i) < 3 {
			continue
		}

		input := i[1]

		if i[2] != "" {
			input = strings.TrimPrefix(i[2], "|")
		}

		id, name, err := db.findTerm(ctx, conn, input)
		if err == nil && db.TermBaseURL != "" {
			// hacky way to make links lowercase if the input was lowercase
			if lowercaseRegexp.Match([]byte{i[1][0]}) {
				name = strings.ToLower(name)
			}

			replace := fmt.Sprintf("[%v](%v%v)", name, db.TermBaseURL, id)
			if len(i) > 2 && i[2] != "" {
				replace = fmt.Sprintf("[%v](%v%v)", i[1], db.TermBaseURL, id)
			}

			s = append(s, i[0], replace)
		} else {
			fmt.Printf("Error fetching term with name or ID `%v`: %v\n", input, err)

			if len(i) > 2 && i[2] != "" {
				s = append(s, i[0], input)
			} else {
				s = append(s, i[0], i[1])
			}
		}
	}

	r := strings.NewReplacer(s...)
	return r.Replace(input)
}

var numberRegex = regexp.MustCompile(`^\d+$`)

func (db *Db) findTerm(ctx context.Context, conn *pgxpool.Conn, in string) (id int, name string, err error) {
	sql := "select id, name from terms where "

	if numberRegex.MatchString(in) {
		sql += "id = $1::int"
	} else {
		sql += "lower(name) = lower($1::text)"
	}

	err = conn.QueryRow(ctx, sql, in).Scan(&id, &name)
	return
}
