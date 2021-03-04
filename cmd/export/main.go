// Package main is a simple utility to export a Berry/Termora database
package main

import (
	"encoding/json"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/termora/berry/commands/static/export"
	"github.com/termora/berry/db"
)

func main() {
	var (
		url            string
		escape, indent bool

		pronouns, explanations bool
	)

	flag.StringVarP(&url, "url", "u", "", "Database URL")
	flag.BoolVarP(&escape, "escape", "e", false, "Escape HTML characters")
	flag.BoolVarP(&indent, "indent", "i", false, "Indent the output (human-readable)")
	flag.BoolVarP(&pronouns, "pronouns", "p", false, "Export pronouns")
	flag.BoolVarP(&explanations, "explanations", "x", false, "Export explanations")
	flag.Parse()
	if url == "" {
		panic("no database url provided")
	}

	db, err := db.Init(url, nil)
	if err != nil {
		panic(err)
	}

	export, err := export.New(db)
	if err != nil {
		panic(err)
	}

	if !pronouns {
		export.Pronouns = nil
	}
	if !explanations {
		export.Explanations = nil
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(escape)
	if indent {
		enc.SetIndent("", "  ")
	}

	if err = enc.Encode(export); err != nil {
		panic(err)
	}
}
