package main

import (
	"net/http"

	"github.com/Starshine113/berry/db"

	"github.com/labstack/echo/v4"
)

func (a *api) add(c echo.Context) (err error) {
	t := new(db.Term)
	if err = c.Bind(t); err != nil {
		return
	}

	t, err = a.db.AddTerm(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, t)
}
