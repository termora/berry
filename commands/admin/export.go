package admin

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/crouter"
	"github.com/bwmarrin/discordgo"
)

const eVersion = 1

type e struct {
	Version    int        `json:"export_version"`
	ExportDate time.Time  `json:"export_date"`
	Terms      []*db.Term `json:"terms"`
}

func (c *commands) export(ctx *crouter.Ctx) (err error) {
	export := e{ExportDate: time.Now().UTC(), Version: eVersion}

	var gz bool
	if strings.Contains(ctx.RawArgs, "-gz") || strings.Contains(ctx.RawArgs, "-gzip") {
		gz = true
	}

	u, err := ctx.Session.UserChannelCreate(ctx.Author.ID)
	if err != nil {
		c.sugar.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?")
		return
	}

	terms, err := c.db.GetTerms(0)
	if err != nil {
		return ctx.CommandError(err)
	}

	export.Terms = terms

	b, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return ctx.CommandError(err)
	}
	fn := fmt.Sprintf("export-%v.json", time.Now().Format("2006-01-02-15-04-05"))

	var buf *bytes.Buffer
	if gz {
		buf = new(bytes.Buffer)
		zw := gzip.NewWriter(buf)
		zw.Name = fn
		_, err = zw.Write(b)
		if err != nil {
			return ctx.CommandError(err)
		}
		err = zw.Close()
		if err != nil {
			return ctx.CommandError(err)
		}
		fn = fn + ".gz"
	} else {
		buf = bytes.NewBuffer(b)
	}

	file := discordgo.File{
		Name:   fn,
		Reader: buf,
	}

	_, err = ctx.Session.ChannelMessageSendComplex(u.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("> Done! Archive of %v terms, invoked by %v at %v.", len(terms), ctx.Author.String(), time.Now().Format(time.RFC3339)),
		Files:   []*discordgo.File{&file},
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeUsers,
			},
		},
	})
	return err
}
