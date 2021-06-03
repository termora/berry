package main

import (
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"github.com/termora/berry/db"
)

func funcMap() template.FuncMap {
	return map[string]interface{}{
		"urlEncode": func(s string) string {
			return url.PathEscape(s)
		},
		"urlDecode": func(s string) string {
			u, err := url.PathUnescape(s)
			if err != nil {
				return ""
			}
			return u
		},
		"timeToDate": func(t time.Time) string {
			return t.Format("Mon January 02 2006")
		},
		"markdownParse": func(s string) template.HTML {
			return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(
				blackfriday.Run(
					[]byte(s),
					blackfriday.WithExtensions(blackfriday.Autolink|blackfriday.Strikethrough|blackfriday.HardLineBreak))))
		},
		"sanitize": func(s string) template.HTML {
			return template.HTML(bluemonday.UGCPolicy().Sanitize(s))
		},
		"resultsNum": func(s []*db.Term) int {
			return len(s)
		},
		"title": strings.Title,
		"pageStyle": func(darkMode string) template.CSS {
			if darkMode == "false" {
				return template.CSS("")
			} else if darkMode == "true" {
				return template.CSS(`body { color: #dcddde; background-color: #36393f; }`)
			} else {
				return template.CSS(`@media (prefers-color-scheme: dark) { body { color: #dcddde; background-color: #36393f; } }`)
			}
		},
	}
}
