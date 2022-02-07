package site

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *site) index(c echo.Context) (err error) {
	tags, err := s.db.Tags()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Render(http.StatusOK, "index.html", (&renderData{
		Conf: s.Config,
		Tags: tags,
	}).parse(c))
}
