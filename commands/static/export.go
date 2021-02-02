package static

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/utils/sendpart"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
)

const eVersion = 1

type e struct {
	Version      int               `json:"export_version"`
	ExportDate   time.Time         `json:"export_date"`
	Terms        []*db.Term        `json:"terms"`
	Explanations []*db.Explanation `json:"explanations"`
}

func (c *Commands) export(ctx *bcr.Context) (err error) {
	export := e{ExportDate: time.Now().UTC(), Version: eVersion}

	var gz bool
	if strings.Contains(ctx.RawArgs, "-gz") || strings.Contains(ctx.RawArgs, "-gzip") {
		gz = true
	}

	u, err := ctx.Session.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		c.Sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?", nil)
		return
	}

	terms, err := c.DB.GetTerms(0)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	export.Terms = terms

	ex, err := c.DB.GetAllExplanations()
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}
	export.Explanations = ex

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

	file := sendpart.File{
		Name:   fn,
		Reader: buf,
	}

	_, err = ctx.Session.SendMessageComplex(u.ID, api.SendMessageData{
		Content:         fmt.Sprintf("> Done! Archive of %v terms, invoked by %v#%v at %v.", len(terms), ctx.Author.Username, ctx.Author.Discriminator, time.Now().Format(time.RFC3339)),
		Files:           []sendpart.File{file},
		AllowedMentions: ctx.Router.DefaultMentions,
	})
	return err
}
