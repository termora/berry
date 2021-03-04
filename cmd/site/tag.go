package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/termora/berry/db"
)

func (s *site) tag(c echo.Context) (err error) {
	var terms []*db.Term
	if c.Param("tag") == "untagged" || c.Param("tag") == "" {
		terms, err = s.db.UntaggedTerms()
	} else {
		terms, err = s.db.TagTerms(c.Param("tag"))
	}
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	data := struct {
		Conf  conf
		Tag   string
		Terms []*db.Term
	}{s.conf, c.Param("tag"), terms}

	return c.Render(http.StatusOK, "terms.html", data)
}
