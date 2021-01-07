package admin

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/misc"
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/utils/sendpart"
)

const eVersion = 1

type e struct {
	Version    int        `json:"export_version"`
	ExportDate time.Time  `json:"export_date"`
	Terms      []*db.Term `json:"terms"`
}

func (c *commands) export(ctx *bcr.Context) (err error) {
	export := e{ExportDate: time.Now().UTC(), Version: eVersion}

	var gz bool
	if strings.Contains(ctx.RawArgs, "-gz") || strings.Contains(ctx.RawArgs, "-gzip") {
		gz = true
	}

	u, err := ctx.Session.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		c.sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?", nil)
		return
	}

	terms, err := c.db.GetTerms(0)
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	export.Terms = terms

	b, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}
	fn := fmt.Sprintf("export-%v.json", time.Now().Format("2006-01-02-15-04-05"))

	var buf *bytes.Buffer
	if gz {
		buf = new(bytes.Buffer)
		zw := gzip.NewWriter(buf)
		zw.Name = fn
		_, err = zw.Write(b)
		if err != nil {
			_, err = ctx.Send(misc.InternalError, nil)
			return err
		}
		err = zw.Close()
		if err != nil {
			_, err = ctx.Send(misc.InternalError, nil)
			return err
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
