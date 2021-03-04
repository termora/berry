package main

import (
	"net/http"

	"github.com/termora/berry/db"

	"github.com/labstack/echo/v4"
)

func (a *api) add(c echo.Context) (err error) {
	t := new(db.Term)
	// parse the request into a Term
	if err = c.Bind(t); err != nil {
		return
	}

	// add the new term to the database
	t, err = a.db.AddTerm(t)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// return the new term
	return c.JSON(http.StatusOK, t)
}
