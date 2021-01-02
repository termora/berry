package main

import (
	"net/http"
	"strconv"

	"github.com/Starshine113/berry/db"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (a *api) list(c echo.Context) (err error) {
	terms, err := a.db.GetTerms(db.FlagListHidden)
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

func (a *api) listCategory(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	terms, err := a.db.GetCategoryTerms(id, db.FlagListHidden)
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

func (a *api) categories(c echo.Context) (err error) {
	categories, err := a.db.GetCategories()
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	if len(categories) == 0 {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, categories)
}

func (a *api) explanations(c echo.Context) (err error) {
	explanations, err := a.db.GetAllExplanations()
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, explanations)
}
