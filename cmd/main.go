// Package main is the main entrypoint for Termora
package main

import (
	"fmt"
	"os"

	"github.com/termora/berry/cmd/api"
	"github.com/termora/berry/cmd/bot"
	"github.com/termora/berry/cmd/export"
	"github.com/termora/berry/cmd/site"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "Termora",
	Usage: "A searchable glossary bot and website",

	Commands: []*cli.Command{
		bot.Command,
		site.Command,
		api.Command,
		export.Command,
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error in command:", err)
		os.Exit(1)
	}
}
