package main

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

//go:embed static/pages/*
var staticPages embed.FS

var pages = map[string]template.HTML{}
var pageMu sync.RWMutex

// limit page names to 64 characters
var pageRegex = regexp.MustCompile(`^\w{1,64}$`)

func (s *site) staticPage(c echo.Context) (err error) {
	name, err := url.PathUnescape(c.Param("page"))
	if err != nil {
		name = ""
	}
	name = strings.ToLower(
		strings.TrimSpace(
			strings.TrimSuffix(name, ".md"),
		),
	)

	data := &renderData{
		Conf: s.conf,
	}

	pageMu.RLock()
	t, ok := pages[name]
	if ok {
		pageMu.RUnlock()
		data.MD = t
		return c.Render(http.StatusOK, "static.html", data.parse(c))
	}
	pageMu.RUnlock()

	if !pageRegex.MatchString(name) {
		return c.Render(http.StatusNotFound, "404.html", data.parse(c))
	}

	b, err := staticPages.ReadFile(path.Join("static/pages/", name+".md"))
	if err != nil {
		return c.Render(http.StatusNotFound, "404.html", data.parse(c))
	}

	text := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(
		blackfriday.Run(b,
			blackfriday.WithExtensions(blackfriday.Autolink|blackfriday.Strikethrough|blackfriday.HardLineBreak))))

	pageMu.Lock()
	pages[name] = text
	pageMu.Unlock()

	data.MD = text

	return c.Render(http.StatusOK, "static.html", data.parse(c))
}
