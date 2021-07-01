package search

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/snowflake/v2"
	"github.com/termora/berry/db"
)

func (c *commands) files(ctx *bcr.Context) (err error) {
	files := []db.File{}
	if ctx.RawArgs == "" {
		files, err = c.DB.Files()
	} else {
		files, err = c.DB.FileName(ctx.RawArgs)
	}

	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	if len(files) == 0 {
		_, err = ctx.Reply("There are no files in the database.")
		return
	}

	s := []string{}

	for _, f := range files {
		if c.Config.Bot.Website == "" {
			s = append(s, fmt.Sprintf("`%v`: %v\n", f.ID, f.Filename))
		} else {
			s = append(s, fmt.Sprintf("`%v`: [%v](%vfile/%v/%v)\n", f.ID, f.Filename, c.Config.Bot.Website, f.ID, f.Filename))
		}
	}

	name := "Files"
	if ctx.RawArgs != "" {
		arg := ctx.RawArgs
		if len(arg) > 100 {
			arg = arg[:100] + "..."
		}

		name += " matching " + arg
	}

	_, err = ctx.PagedEmbed(bcr.StringPaginator(name, db.EmbedColour, s, 10), false)
	return
}

func (c *commands) file(ctx *bcr.Context) (err error) {
	i, err := strconv.ParseUint(ctx.RawArgs, 0, 0)
	if err == nil {
		f, err := c.DB.File(snowflake.ID(i))
		if err != nil {
			_, err = ctx.Reply("No file with that ID found.")
			return err
		}

		e := discord.Embed{
			Title:       f.Filename,
			Description: fmt.Sprintf("[Link](%v)", f.URL()),
			Image: &discord.EmbedImage{
				URL: f.URL(),
			},
			Color: db.EmbedColour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("ID: %v | Added", f.ID),
			},
			Timestamp: discord.NewTimestamp(c.DB.Time(f.ID)),
		}

		if f.Description != "" {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "Description",
				Value: f.Description,
			})
		}

		if f.Source != "" {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "Source",
				Value: f.Source,
			})
		}

		_, err = ctx.Send("", &e)
		return err
	}

	files, err := c.DB.FileName(ctx.RawArgs)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	if len(files) == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "No files with that name found.")
		return
	}

	if len(files) > 1 {
		_, err = ctx.Replyc(bcr.ColourRed, "More than one file with that name found. Try `%vfiles %v`?", ctx.Prefix, ctx.RawArgs)
		return
	}

	f := files[0]

	e := discord.Embed{
		Title:       f.Filename,
		Description: fmt.Sprintf("[Link](%v)", f.URL()),
		Image: &discord.EmbedImage{
			URL: f.URL(),
		},
		Color: db.EmbedColour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Added", f.ID),
		},
		Timestamp: discord.NewTimestamp(c.DB.Time(f.ID)),
	}

	if f.Description != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Description",
			Value: f.Description,
		})
	}

	if f.Source != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Source",
			Value: f.Source,
		})
	}

	_, err = ctx.Send("", &e)
	return err
}
