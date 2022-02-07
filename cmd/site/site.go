package site

import (
	"context"
	"embed"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/termora/berry/db"
	"github.com/termora/berry/db/search/typesense"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

//go:embed templates/*
var tmpls embed.FS

//go:embed static
var staticFS embed.FS

func mustSub(f fs.FS, path string) fs.FS {
	sub, err := fs.Sub(f, path)
	if err != nil {
		panic(err)
	}
	return sub
}

var Command = &cli.Command{
	Name:    "site",
	Aliases: []string{"web"},
	Usage:   "Run the website",
	Action:  run,
}

type site struct {
	db    *db.DB
	conf  conf
	sugar *zap.SugaredLogger
}

// T ...
type T struct {
	templates *template.Template
}

// Render ...
func (t *T) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type conf struct {
	DatabaseURL string `yaml:"database_url"`
	Port        string

	SiteName string `yaml:"site_name"`
	BaseURL  string `yaml:"base_url"`
	Invite   string `yaml:"invite_url"`
	Git      string
	Contact  bool
	// Optional description shown in embeds, when not linking to a term page
	Description string

	Plausible struct {
		Domain string
		URL    string
	}

	Typesense struct {
		URL string
		Key string
	}
}

type renderData struct {
	Conf  conf
	Path  string
	Dark  string
	Tag   string
	Tags  []string
	Term  *db.Term
	Terms []*db.Term
	Query template.HTML
	// Parsed markdown text for about pages
	MD template.HTML
}

func (r *renderData) parse(c echo.Context) renderData {
	r.Path = c.Request().URL.Path

	if cookie, err := c.Request().Cookie("dark"); err == nil {
		r.Dark = cookie.Value
	} else {
		r.Dark = ""
	}

	return *r
}

func run(ctx *cli.Context) error {
	t := &T{
		templates: template.Must(template.New("").
			Funcs(sprig.FuncMap()).
			Funcs(funcMap()).
			ParseFS(tmpls, "templates/*.html")),
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	sugar := logger.Sugar()

	var c conf

	configFile, err := ioutil.ReadFile("config.site.yaml")
	if err != nil {
		sugar.Fatal(err)
	}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		sugar.Fatalf("Error loading configuration file: %v", err)
	}
	sugar.Info("Loaded configuration file.")

	d, err := db.Init(c.DatabaseURL, sugar)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	d.TermBaseURL = "/term/"
	sugar.Info("Connected to database.")

	// Typesense requires a bot running to sync terms
	if c.Typesense.URL != "" && c.Typesense.Key != "" {
		d.Searcher, err = typesense.New(c.Typesense.URL, c.Typesense.Key, d.Pool, sugar.Debugf)
		if err != nil {
			sugar.Fatalf("Couldn't connect to Typesense: %v", err)
		}
		sugar.Info("Connected to Typesense server")
	}

	s := site{db: d, conf: c, sugar: sugar}

	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Logger())

	e.GET("/static/*", echo.WrapHandler(
		http.StripPrefix("/static/", http.FileServer(http.FS(mustSub(staticFS, "static")))),
	))

	e.GET("/dark", setDarkPreferences)

	e.GET("/", s.index)
	e.GET("/term/:term", s.term)
	e.GET("/tag/:tag", s.tag)
	e.GET("/search", s.search)
	e.GET("/file/:id/:filename", s.file)
	e.GET("/about/:page", s.staticPage)

	e.GET("/robots.txt", func(c echo.Context) error {
		return c.String(http.StatusOK, `User-agent: *
Disallow: /file
Disallow: /search
Disallow: /static`)
	})

	// get port
	port := c.Port

	if port == "" {
		port = "1300"
	} else {
		port = strings.TrimPrefix(port, ":")
	}

	go func() {
		if err := e.Start(":" + port); err != nil {
			sugar.Info("Shutting down server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cctx, cancel := context.WithTimeout(ctx.Context, 10*time.Second)
	defer cancel()
	if err := e.Shutdown(cctx); err != nil {
		sugar.Fatal(err)
	}
	return err
}
