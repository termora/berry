package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	flags := search.FlagListHidden
	if s := chi.URLParam(r, "flags"); s != "" {
		f, _ := strconv.Atoi(s)
		flags = search.TermFlag(f)
	}

	terms, err := s.db.GetTerms(flags)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.log.Errorf("Error getting terms: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, terms)
}

func (s *Server) listCategory(w http.ResponseWriter, r *http.Request) {
	// parse the ID
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	// get all terms from that category
	terms, err := s.db.GetCategoryTerms(id, search.FlagListHidden)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.log.Errorf("Error getting terms in category %v: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, terms)
}

func (s *Server) categories(w http.ResponseWriter, r *http.Request) {
	// get all categories
	categories, err := s.db.GetCategories()
	if err != nil {
		s.log.Errorf("Error getting categories: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, categories)
}

func (s *Server) explanations(w http.ResponseWriter, r *http.Request) {
	// get all explanations
	explanations, err := s.db.GetAllExplanations()
	if err != nil {
		s.log.Errorf("Error getting explanations: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, explanations)
}

func (s *Server) pronouns(w http.ResponseWriter, r *http.Request) {
	pronouns, err := s.db.Pronouns(db.AlphabeticPronounOrder)
	if err != nil {
		s.log.Errorf("Error getting pronouns: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, pronouns)
}

func (s *Server) tags(w http.ResponseWriter, r *http.Request) {
	tags, err := s.db.Tags()
	if err != nil {
		s.log.Errorf("Error getting tags: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, tags)
}
