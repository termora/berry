// Package typesense implements search methods with a Typesense search server.
package typesense

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/termora/berry/db/search"
	"github.com/typesense/typesense-go/typesense"
)

// New returns a new Searcher
func New(dsn, apiKey string, pg *pgxpool.Pool, debugFunc func(string, ...interface{})) (search.Searcher, error) {
	if debugFunc == nil {
		debugFunc = func(string, ...interface{}) {
			return
		}
	}

	c := typesense.NewClient(
		typesense.WithServer(dsn),
		typesense.WithAPIKey(apiKey),
		typesense.WithConnectionTimeout(10*time.Second),
	)

	_, err := c.Health(10 * time.Second)
	if err != nil {
		return nil, err
	}

	return &Client{
		ts:    c,
		pg:    pg,
		Debug: debugFunc,
	}, nil
}

var _ search.Searcher = (*Client)(nil)

// Client ...
type Client struct {
	ts *typesense.Client
	pg *pgxpool.Pool

	Debug func(string, ...interface{})
}
