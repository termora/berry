package static

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/static/export"
	"github.com/termora/berry/db"
)

const eVersion = 3

type e struct {
	Version      int               `json:"export_version"`
	ExportDate   time.Time         `json:"export_date"`
	Categories   []*db.Category    `json:"categories"`
	Terms        []*db.Term        `json:"terms"`
	Tags         []string          `json:"tags"`
	Explanations []*db.Explanation `json:"explanations,omitempty"`
	Pronouns     []*db.PronounSet  `json:"pronouns,omitempty"`
}

func (c *Commands) export(ctx *bcr.Context) (err error) {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	var gz bool
	fs.BoolVarP(&gz, "compress", "x", false, "Compress the output with gzip")

	u, err := ctx.Session.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		c.Sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?", nil)
		return
	}

	export, err := export.New(c.DB)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	b, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	fn := fmt.Sprintf("export-%v.json", time.Now().Format("2006-01-02-15-04-05"))

	var buf *bytes.Buffer
	if gz {
		buf = new(bytes.Buffer)
		zw := gzip.NewWriter(buf)
		zw.Name = fn
		_, err = zw.Write(b)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
		err = zw.Close()
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
		fn = fn + ".gz"
	} else {
		buf = bytes.NewBuffer(b)
	}

	_, err = ctx.NewMessage().Channel(u.ID).TogglePermCheck().Content(
		fmt.Sprintf(
			"> Done! Archive of %v terms, %v explanations, and %v pronoun sets.",
			len(export.Terms), len(export.Explanations), len(export.Pronouns),
		),
	).AddFile(fn, buf).Send()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Send("✅ Check your DMs!", nil)
	return err
}
