package main

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/termora/berry/db"
)

func (s *site) tag(c echo.Context) (err error) {
	tag, err := url.PathUnescape(c.Param("tag"))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var terms []*db.Term
	if tag == "untagged" || tag == "" {
		terms, err = s.db.UntaggedTerms()
	} else {
		terms, err = s.db.TagTerms(tag)
	}
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "terms.html", renderData{
		Conf:  s.conf,
		Tag:   tag,
		Terms: terms,
	})
}
