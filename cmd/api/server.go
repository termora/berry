package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/termora/berry/db"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type api struct {
	db    *db.Db
	conf  conf
	sugar *zap.SugaredLogger
}

type conf struct {
	DatabaseURL string `yaml:"database_url"`
	Port        string `yaml:"port"`
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// read config
	var c conf

	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &c)
	sugar.Info("Loaded configuration file.")

	// connect to the database
	d, err := db.Init(c.DatabaseURL, sugar)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	sugar.Info("Connected to database.")

	// create the API
	r := api{db: d, conf: c, sugar: sugar}

	// create the router
	e := echo.New()
	// add logging
	e.Use(middleware.Logger())

	// add the routes
	e.GET("/v1/search/:term", r.search)
	e.GET("/v1/term/:id", r.term)
	e.GET("/v1/list", r.list)
	e.GET("/v1/list/:id", r.listCategory)
	e.GET("/v1/categories", r.categories)
	e.GET("/v1/explanations", r.explanations)

	e.GET("/robots.txt", func(c echo.Context) error {
		return c.String(http.StatusOK, `User-agent: *
Disallow: /`)
	})

	// get port
	port := c.Port
	strings.TrimPrefix(port, ":")
	if port == "" {
		port = "1300"
	}

	// run the server
	go func() {
		if err := e.Start(":" + c.Port); err != nil {
			sugar.Info("Shutting down server")
		}
	}()

	// gracefully catch interrupts
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		sugar.Fatal(err)
	}
}
