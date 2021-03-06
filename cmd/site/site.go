package main

import (
	"context"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/termora/berry/db"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type site struct {
	db    *db.Db
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

	Plausible struct {
		Domain string
		URL    string
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

func main() {
	t := &T{
		templates: template.Must(template.New("").
			Funcs(sprig.FuncMap()).
			Funcs(funcMap()).
			ParseGlob("templates/*.html")),
	}

	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	var c conf

	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
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
	sugar.Info("Connected to database.")

	s := site{db: d, conf: c, sugar: sugar}

	e := echo.New()
	e.Renderer = t
	e.Use(middleware.Logger())
	e.Static("/static", "static")

	e.GET("/dark", setDarkPreferences)

	e.GET("/", s.index)
	e.GET("/term/:term", s.term)
	e.GET("/tag/:tag", s.tag)
	e.GET("/search", s.search)
	e.GET("/file/:id/:filename", s.file)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		sugar.Fatal(err)
	}
}
