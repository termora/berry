// Package typesense implements search methods with a Typesense search server.
package typesense

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/termora/berry/db/search"
	"github.com/termora/tsclient"
)

// New returns a new Searcher
func New(dsn, apiKey string, pg *pgxpool.Pool) (search.Searcher, error) {

	c, err := tsclient.New(dsn, apiKey)
	if err != nil {
		return nil, err
	}

	return &Client{
		ts: c,
		pg: pg,
	}, nil
}

var _ search.Searcher = (*Client)(nil)

// Client ...
type Client struct {
	ts *tsclient.Client
	pg *pgxpool.Pool
}
