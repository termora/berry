package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/starshine-sys/berry/db"
)

type searchData struct {
	Conf  conf
	Terms []*db.Term
	Query template.HTML
}

func (s *site) search(c echo.Context) (err error) {
	q := template.HTML(bluemonday.UGCPolicy().Sanitize(c.QueryParam("q")))
	terms, err := s.db.Search(c.QueryParam("q"), 0)
	if err != nil {
		return c.Render(http.StatusNotFound, "noQuery.html", searchData{Conf: s.conf, Query: q})
	}

	if len(terms) == 0 {
		return c.Render(http.StatusNotFound, "noQuery.html", searchData{Conf: s.conf, Query: q})
	}

	return c.Render(http.StatusOK, "results.html",
		searchData{
			Conf:  s.conf,
			Terms: terms,
			Query: q,
		},
	)
}
