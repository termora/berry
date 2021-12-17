package main

import (
	"net/http"
	"strconv"

	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (a *api) list(c echo.Context) (err error) {
	flags := search.FlagListHidden
	if c.QueryParam("flags") != "" {
		f, _ := strconv.Atoi(c.QueryParam("flags"))
		flags = search.TermFlag(f)
	}

	terms, err := a.db.GetTerms(flags)
	if err != nil {
		// if no rows were returned, return no content
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		// otherwise, internal server error
		a.sugar.Errorf("Error getting terms: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// if no rows were returned, return no content
	if len(terms) == 0 {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, terms)
}

func (a *api) listCategory(c echo.Context) (err error) {
	// parse the ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// get all terms from that category
	terms, err := a.db.GetCategoryTerms(id, search.FlagListHidden)
	if err != nil {
		// if no rows were returned, return no content
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		a.sugar.Errorf("Error getting terms: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// if no rows were returned, return no content
	if len(terms) == 0 {
		return c.NoContent(http.StatusNoContent)
	}
	return c.JSON(http.StatusOK, terms)
}

func (a *api) categories(c echo.Context) (err error) {
	// get all categories
	categories, err := a.db.GetCategories()
	if err != nil {
		// if no rows were returned, return no content
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		a.sugar.Errorf("Error getting categories: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if len(categories) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	// if no rows were returned, return no content
	return c.JSON(http.StatusOK, categories)
}

func (a *api) explanations(c echo.Context) (err error) {
	// get all explanations
	explanations, err := a.db.GetAllExplanations()
	if err != nil {
		// if no rows were returned, return no content
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		a.sugar.Errorf("Error getting explanations: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, explanations)
}

func (a *api) pronouns(c echo.Context) (err error) {
	pronouns, err := a.db.Pronouns(db.AlphabeticPronounOrder)
	if err != nil {
		// if no rows were returned, return no content
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		// otherwise, internal server error
		a.sugar.Errorf("Error getting pronouns: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// if no rows were returned, return no content
	if len(pronouns) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, pronouns)
}

func (a *api) tags(c echo.Context) (err error) {
	tags, err := a.db.Tags()
	if err != nil {
		// if no rows were returned, return no content
		if errors.Cause(err) == pgx.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		// otherwise, internal server error
		a.sugar.Errorf("Error getting tags: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// if no rows were returned, return no content
	if len(tags) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, tags)
}
