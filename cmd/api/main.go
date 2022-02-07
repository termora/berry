package api

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search/typesense"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var Command = &cli.Command{
	Name:   "api",
	Usage:  "Run the API",
	Action: run,
}

type Server struct {
	db   *db.DB
	conf conf
	log  *zap.SugaredLogger
}

type conf struct {
	DatabaseURL  string `yaml:"database_url"`
	TypesenseURL string `yaml:"typesense_url"`
	TypesenseKey string `yaml:"typesense_key"`
	Port         string `yaml:"port"`
}

func run(*cli.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log := logger.Sugar()

	// read config
	var c conf

	configFile, err := ioutil.ReadFile("config.api.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Loaded configuration file.")

	s := &Server{conf: c, log: log}

	// connect to the database
	s.db, err = db.Init(c.DatabaseURL, log)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Info("Connected to database")

	if s.conf.TypesenseURL != "" && s.conf.TypesenseKey != "" {
		s.db.Searcher, err = typesense.New(s.conf.TypesenseURL, s.conf.TypesenseKey, s.db.Pool, s.log.Debugf)
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
		r.Get(`/id/{id:\d+}`, s.term)

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
	port := c.Port
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
