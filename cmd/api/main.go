package api

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/termora/berry/common"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search/typesense"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:   "api",
	Usage:  "Run the API",
	Action: run,
}

type Server struct {
	db *db.DB
}

func run(*cli.Context) (err error) {
	// read config
	c := common.ReadConfig()

	log.Info("Loaded configuration file.")

	s := &Server{}

	// connect to the database
	s.db, err = db.Init(c.Core.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Info("Connected to database")

	if c.Core.TypesenseURL != "" && c.Core.TypesenseKey != "" {
		s.db.Searcher, err = typesense.New(c.Core.TypesenseURL, c.Core.TypesenseKey, s.db.Pool)
		if err != nil {
			log.Fatalf("Error connecting to Typesense: %v", err)
		}
		log.Info("Connected to Typesense")
	}

	mx := chi.NewMux()
	mx.Use(middleware.Recoverer)
	mx.Use(middleware.RedirectSlashes)
	mx.Use(middleware.CleanPath)

	mx.Route("/v1", func(r chi.Router) {
		r.Get("/search/{term}", s.search)
		r.Get(`/term/{id:\d+}`, s.term)

		r.Get("/list", s.list)
		r.Get(`/list/{id:\d+}`, s.listCategory)

		r.Get("/categories", s.categories)
		r.Get("/explanations", s.explanations)
		r.Get("/explanations", s.explanations)
		r.Get("/tags", s.tags)
		r.Get("/pronouns", s.pronouns)
	})

	mx.Get("/robots.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`User-agent: *
Disallow: /`))
	})

	// get port
	port := c.API.Port
	if port == "" {
		port = ":1300"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	e := make(chan error)

	go func() {
		e <- http.ListenAndServe(port, mx)
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	log.Infof("Listening on %v!", port)
	log.Infof("Press Ctrl-C or send an interrupt signal to stop.")

	select {
	case <-sc:
		log.Infof("Interrupt signal received. Shutting down...")
		s.db.Close()
	case err := <-e:
		log.Errorf("Error serving API: %v", err)
	}
	return nil
}
