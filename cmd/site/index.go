package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *site) index(c echo.Context) (err error) {
	tags, err := s.db.Tags()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	data := struct {
		Conf conf
		Tags []string
	}{s.conf, tags}

	return c.Render(http.StatusOK, "index.html", data)
}
