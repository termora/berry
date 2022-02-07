package site

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/snowflake/v2"
)

func (s *site) file(c echo.Context) (err error) {
	if c.Param("id") == "" {
		return c.NoContent(http.StatusNotFound)
	}
	i, err := strconv.ParseUint(c.Param("id"), 0, 0)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	f, err := s.db.File(snowflake.ID(i))
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNotFound)
		}

		s.sugar.Errorf("Error getting file %v: %v", i, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, f.ContentType, f.Data)
}
