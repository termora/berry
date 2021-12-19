package admin

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) upload(ctx *bcr.Context) (err error) {
	if len(ctx.Message.Attachments) == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "No files attached.")
		return
	}

	a := ctx.Message.Attachments[0]

	if a.Size > 1*1024*1024 {
		_, err = ctx.Replyc(bcr.ColourRed, "The attachment is too big (%v (%v bytes) > 1 MB)", humanize.Bytes(a.Size), humanize.Comma(int64(a.Size)))
		return
	}

	contentType := contentType(a.Filename)

	if contentType == "unknown" {
		_, err = ctx.Replyc(bcr.ColourRed, "The attachment you gave isn't an image.")
		return
	}

	resp, err := http.Get(a.URL)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}
	defer resp.Body.Close()

	ctx.State.Typing(ctx.Channel.ID)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	f, err := bot.DB.AddFile(a.Filename, contentType, data)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	link := ""
	if f.URL() != "" {
		link = fmt.Sprintf(" ([link](%v))", f.URL())
	}

	_, err = ctx.Reply("File added with ID %v!%v", f.ID, link)
	return
}

func contentType(filename string) string {
	if bcr.HasAnySuffix(filename, ".jpg", ".jpeg") {
		return "image/jpeg"
	}

	if strings.HasSuffix(filename, ".png") {
		return "image/png"
	}

	if strings.HasSuffix(filename, ".gif") {
		return "image/gif"
	}

	if strings.HasSuffix(filename, ".webp") {
		return "image/webp"
	}

	return "unknown"
}
