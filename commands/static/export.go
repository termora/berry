package static

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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

	u, err := ctx.State.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		c.Sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?")
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

	if ctx.Channel.ID != u.ID {
		_, err = ctx.Send("✅ Check your DMs!")
	}

	return err
}

func (c *Commands) exportCSV(ctx *bcr.Context) (err error) {
	terms, err := c.DB.GetTerms(0)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	var b bytes.Buffer

	w := csv.NewWriter(&b)
	w.Write([]string{"ID", "Term", "Description", "Coined by", "Tags"})
	for _, t := range terms {
		w.Write([]string{
			strconv.Itoa(t.ID),
			strings.Join(append([]string{t.Name}, t.Aliases...), ", "),
			t.Description,
			t.Source,
			strings.Join(t.DisplayTags, ", "),
		})
	}

	w.Flush()

	err = w.Error()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.NewMessage().Content("Here you go!").AddFile("terms.csv", &b).Send()
	return
}

func (c *Commands) exportXLSX(ctx *bcr.Context) (err error) {
	terms, err := c.DB.GetTerms(0)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	f := excelize.NewFile()

	// set header
	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Term")
	f.SetCellValue("Sheet1", "C1", "Description")
	f.SetCellValue("Sheet1", "D1", "Coined by")
	f.SetCellValue("Sheet1", "E1", "Tags")

	for i, t := range terms {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i+2), t.ID)
		f.SetCellValue(
			"Sheet1", fmt.Sprintf("B%v", i+2),
			strings.Join(append([]string{t.Name}, t.Aliases...), ", "),
		)
		f.SetCellValue(
			"Sheet1", fmt.Sprintf("C%v", i+2),
			t.Description,
		)
		f.SetCellValue(
			"Sheet1", fmt.Sprintf("D%v", i+2),
			t.Source,
		)
		f.SetCellValue(
			"Sheet1", fmt.Sprintf("E%v", i+2),
			strings.Join(t.DisplayTags, ", "),
		)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.NewMessage().Content("Here you go!").AddFile("terms.xlsx", buf).Send()
	return
}
