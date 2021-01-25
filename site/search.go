package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starshine-sys/berry/db"
)

type searchData struct {
	Conf  conf
	Terms []*db.Term
}

func (s *site) search(c echo.Context) (err error) {
	terms, err := s.db.Search(c.QueryParam("q"), 0)
	if err != nil {
		return c.Render(http.StatusNotFound, "noQuery.html", indexData{Conf: s.conf})
	}

	if len(terms) == 0 {
		return c.Render(http.StatusNotFound, "noQuery.html", indexData{Conf: s.conf})
	}

	return c.Render(http.StatusOK, "results.html", searchData{Conf: s.conf, Terms: terms})
}
