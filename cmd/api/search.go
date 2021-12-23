package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "term")

	terms, err := s.db.Search(query, 0, []string{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(terms) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	render.JSON(w, r, terms)
}

func (s *Server) term(w http.ResponseWriter, r *http.Request) {
	// parse the id
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	term, err := s.db.GetTerm(id)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.log.Errorf("Error getting term ID %v: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, term)
}
