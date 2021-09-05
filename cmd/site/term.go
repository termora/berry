package main

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/termora/berry/db"
)

var numberRegex = regexp.MustCompile(`^\d+$`)

func (s *site) term(c echo.Context) (err error) {
	var t *db.Term

	name := c.Param("term")
	name, err = url.PathUnescape(name)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if numberRegex.MatchString(name) {
		id, _ := strconv.Atoi(name)

		t, err = s.db.GetTerm(id)
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}
	} else {
		terms, err := s.db.GetTerms(0)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		for _, i := range terms {
			if strings.EqualFold(i.Name, name) {
				t = i
				break
			}
		}
	}

	if t == nil {
		return c.NoContent(http.StatusNotFound)
	}

	t.Description = s.db.LinkTerms(t.Description)
	t.Note = s.db.LinkTerms(t.Note)
	if t.Disputed() {
		t.Note = strings.TrimSpace(t.Note + "\n\n" + db.DisputedText)
	}

	t.ContentWarnings = s.db.LinkTerms(t.ContentWarnings)

	return c.Render(http.StatusOK, "term.html", (&renderData{
		Conf: s.conf,
		Term: t,
	}).parse(c))
}
