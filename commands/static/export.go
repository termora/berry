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
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/static/export"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
)

func (bot *Bot) export(ctx *bcr.Context) (err error) {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	var gz bool
	fs.BoolVarP(&gz, "compress", "x", false, "Compress the output with gzip")

	u, err := ctx.State.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		log.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?")
		return
	}

	export, err := export.New(bot.DB)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	b, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}
	fn := fmt.Sprintf("export-%v.json", time.Now().Format("2006-01-02-15-04-05"))

	var buf *bytes.Buffer
	if gz {
		buf = new(bytes.Buffer)
		zw := gzip.NewWriter(buf)
		zw.Name = fn
		_, err = zw.Write(b)
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		}
		err = zw.Close()
		if err != nil {
			return bot.DB.InternalError(ctx, err)
		}
		fn = fn + ".gz"
	} else {
		buf = bytes.NewBuffer(b)
	}

	data := api.SendMessageData{
		Content: fmt.Sprintf(
			"> Done! Archive of %v terms, %v explanations, and %v pronoun sets.",
			len(export.Terms), len(export.Explanations), len(export.Pronouns),
		),
		Files: []sendpart.File{{
			Name:   fn,
			Reader: buf,
		}},
	}

	if bot.Config.Bot.LicenseLink != "" {
		data.Embeds = append(data.Embeds, discord.Embed{
			Description: fmt.Sprintf("Make sure to follow the [license](%v).", bot.Config.Bot.LicenseLink),
			Color:       db.EmbedColour,
		})
	}

	_, err = ctx.State.SendMessageComplex(u.ID, data)
	if err != nil {
		return err
	}

	if ctx.Channel.ID != u.ID {
		_, err = ctx.Send("✅ Check your DMs!")
	}

	return err
}

func (bot *Bot) exportCSV(ctx *bcr.Context) (err error) {
	terms, err := bot.DB.GetTerms(0)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	u, err := ctx.State.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		log.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?")
		return
	}

	var b bytes.Buffer

	w := csv.NewWriter(&b)
	_ = w.Write([]string{"ID", "Term", "Description", "Coined by", "Tags"})
	for _, t := range terms {
		_ = w.Write([]string{
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
		return bot.DB.InternalError(ctx, err)
	}

	data := api.SendMessageData{
		Content: "Here you go!",
		Files: []sendpart.File{{
			Name:   "terms.csv",
			Reader: &b,
		}},
	}

	if bot.Config.Bot.LicenseLink != "" {
		data.Embeds = append(data.Embeds, discord.Embed{
			Description: fmt.Sprintf("Make sure to follow the [license](%v).", bot.Config.Bot.LicenseLink),
			Color:       db.EmbedColour,
		})
	}

	_, err = ctx.State.SendMessageComplex(u.ID, data)
	if err != nil {
		return err
	}

	if ctx.Channel.ID != u.ID {
		_, err = ctx.Send("✅ Check your DMs!")
	}

	return err
}

func (bot *Bot) exportXLSX(ctx *bcr.Context) (err error) {
	terms, err := bot.DB.GetTerms(0)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	u, err := ctx.State.CreatePrivateChannel(ctx.Author.ID)
	if err != nil {
		log.Errorf("Error creating user channel for %v: %v", ctx.Author.ID, err)
		_, err = ctx.Send("There was an error opening a DM channel. Are you sure your DMs are open?")
		return
	}

	f := excelize.NewFile()

	// set header
	_ = f.SetCellValue("Sheet1", "A1", "ID")
	_ = f.SetCellValue("Sheet1", "B1", "Term")
	_ = f.SetCellValue("Sheet1", "C1", "Description")
	_ = f.SetCellValue("Sheet1", "D1", "Coined by")
	_ = f.SetCellValue("Sheet1", "E1", "Tags")

	for i, t := range terms {
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i+2), t.ID)
		_ = f.SetCellValue(
			"Sheet1", fmt.Sprintf("B%v", i+2),
			strings.Join(append([]string{t.Name}, t.Aliases...), ", "),
		)
		_ = f.SetCellValue(
			"Sheet1", fmt.Sprintf("C%v", i+2),
			t.Description,
		)
		_ = f.SetCellValue(
			"Sheet1", fmt.Sprintf("D%v", i+2),
			t.Source,
		)
		_ = f.SetCellValue(
			"Sheet1", fmt.Sprintf("E%v", i+2),
			strings.Join(t.DisplayTags, ", "),
		)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	data := api.SendMessageData{
		Content: "Here you go!",
		Files: []sendpart.File{{
			Name:   "terms.xlsx",
			Reader: buf,
		}},
	}

	if bot.Config.Bot.LicenseLink != "" {
		data.Embeds = append(data.Embeds, discord.Embed{
			Description: fmt.Sprintf("Make sure to follow the [license](%v).", bot.Config.Bot.LicenseLink),
			Color:       db.EmbedColour,
		})
	}

	_, err = ctx.State.SendMessageComplex(u.ID, data)
	if err != nil {
		return err
	}

	if ctx.Channel.ID != u.ID {
		_, err = ctx.Send("✅ Check your DMs!")
	}

	return err
}
