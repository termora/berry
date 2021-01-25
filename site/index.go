package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/starshine-sys/berry/db"
)

type indexData struct {
	Conf  conf
	Terms []*db.Term
}

func (s *site) index(c echo.Context) (err error) {
	terms, err := s.db.GetTerms(db.FlagListHidden)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "index.html", indexData{Conf: s.conf, Terms: terms})
}
