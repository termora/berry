package export

import (
	"encoding/json"
	"os"

	"emperror.dev/errors"
	"github.com/termora/berry/commands/static/export"
	"github.com/termora/berry/db"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var Command = &cli.Command{
	Name:   "export",
	Usage:  "Export the database",
	Action: run,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "Database URL",
		},
		&cli.BoolFlag{
			Name:    "escape",
			Aliases: []string{"e"},
			Value:   false,
			Usage:   "Escape HTML characters",
		},
		&cli.BoolFlag{
			Name:    "indent",
			Aliases: []string{"i"},
			Value:   false,
			Usage:   "Indent the output (human-readable)",
		},
		&cli.BoolFlag{
			Name:    "pronouns",
			Aliases: []string{"p"},
			Value:   false,
			Usage:   "Export pronouns",
		},
		&cli.BoolFlag{
			Name:    "explanations",
			Aliases: []string{"x"},
			Value:   false,
			Usage:   "Export explanations",
		},
	},
}

func run(c *cli.Context) error {
	var (
		url    = c.String("url")
		escape = c.Bool("escape")
		indent = c.Bool("indent")

		pronouns     = c.Bool("pronouns")
		explanations = c.Bool("explanations")
	)

	if url == "" {
		return errors.Sentinel("no database url provided")
	}

	db, err := db.Init(url, zap.S())
	if err != nil {
		return err
	}

	export, err := export.New(db)
	if err != nil {
		return err
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
		return err
	}
	return nil
}
