package main

import (
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (a *api) list(c echo.Context) (err error) {
	terms, err := a.db.GetTerms(0)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	if len(terms) == 0 {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, terms)
}
