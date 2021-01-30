package admin

import (
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
)

func (c *Admin) editTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(3); err != nil {
		_, err = ctx.Send("Not enough arguments. Valid subcommands are: `title`, `desc`, `source`, `aliases`.", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Args[1])
	if err != nil {
		_, err = ctx.Sendf("Could not parse ID:\n```%v```", err)
		return
	}
	t, err := c.db.GetTerm(id)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("No term with that ID found.", nil)
			return
		}

		return c.db.InternalError(ctx, err)
	}

	// these should probably be actual subcommands but then we'd have to duplicate the code above 5 times
	switch ctx.Args[0] {
	case "name", "title":
		return c.editTermTitle(ctx, t)
	case "desc", "description":
		return c.editTermDesc(ctx, t)
	case "source":
		return c.editTermSource(ctx, t)
	case "image":
		return c.editTermImage(ctx, t)
	case "aliases":
		return c.editTermAliases(ctx, t)
	}

	_, err = ctx.Send("Invalid subcommand supplied.\nValid subcommands are: `title`, `desc`, `source`, `aliases`.", nil)
	return
}

func (c *Admin) editTermTitle(ctx *bcr.Context, t *db.Term) (err error) {
	title := strings.Join(ctx.Args[2:], " ")
	if len(title) > 200 {
		_, err = ctx.Sendf("Title too long (%v > 200).", len(title))
		return
	}

	err = c.db.UpdateTitle(t.ID, title)
	if err != nil {
		_, err = ctx.Sendf("Error updating title: ```%v```", err)
		return
	}

	_, err = ctx.Send("Title updated!", nil)
	return
}

func (c *Admin) editTermDesc(ctx *bcr.Context, t *db.Term) (err error) {
	desc := strings.Join(ctx.Args[2:], " ")
	if len(desc) > 1800 {
		_, err = ctx.Sendf("Description too long (%v > 1800).", len(desc))
		return
	}

	err = c.db.UpdateDesc(t.ID, desc)
	if err != nil {
		_, err = ctx.Sendf("Error updating description: ```%v```", err)
		return
	}

	_, err = ctx.Send("Description updated!", nil)
	return
}

func (c *Admin) editTermSource(ctx *bcr.Context, t *db.Term) (err error) {
	source := strings.Join(ctx.Args[2:], " ")
	if len(source) > 200 {
		_, err = ctx.Sendf("Source too long (%v > 200).", len(source))
		return
	}

	err = c.db.UpdateSource(t.ID, source)
	if err != nil {
		_, err = ctx.Sendf("Error updating source: ```%v```", err)
		return
	}

	_, err = ctx.Send("Source updated!", nil)
	return
}

func (c *Admin) editTermAliases(ctx *bcr.Context, t *db.Term) (err error) {
	var aliases []string
	if ctx.Args[2] != "clear" {
		aliases = strings.Split(strings.Join(ctx.Args[2:], " "), "\n")
	}

	if len(strings.Join(aliases, ", ")) > 1000 {
		_, err = ctx.Sendf("Total length of aliases too long (%v > 1000)", len(strings.Join(aliases, ", ")))
		return
	}

	err = c.db.UpdateAliases(t.ID, aliases)
	if err != nil {
		_, err = ctx.Sendf("Error updating aliases: ```%v```", err)
		return
	}

	_, err = ctx.Send("Aliases updated!", nil)
	return
}

func (c *Admin) editTermImage(ctx *bcr.Context, t *db.Term) (err error) {
	img := strings.Join(ctx.Args[2:], " ")
	if img == "clear" {
		img = ""
	}

	err = c.db.UpdateImage(t.ID, img)
	if err != nil {
		_, err = ctx.Sendf("Error updating image: ```%v```", err)
		return
	}

	_, err = ctx.Send("Image updated!", nil)
	return
}
