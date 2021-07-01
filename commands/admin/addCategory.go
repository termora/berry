package admin

import (
	"github.com/starshine-sys/bcr"
)

func (c *Admin) addCategory(ctx *bcr.Context) (err error) {
	// if there's no arguments return
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to give a category name.", nil)
		return err
	}

	// check if a category with that name exists
	var e bool

	con, cancel := c.DB.Context()
	defer cancel()

	err = c.DB.Pool.QueryRow(con, "select exists (select from categories where lower(name) = lower($1))", ctx.RawArgs).Scan(&e)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	// if so, return
	if e {
		_, err = ctx.Send(":x :A category with that name already exists.", nil)
		return err
	}

	// add the category
	var id int

	con, cancel = c.DB.Context()
	defer cancel()

	err = c.DB.Pool.QueryRow(con, "insert into public.categories (name) values ($1) returning id", ctx.RawArgs).Scan(&id)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Added category `%v` with ID %v.", ctx.RawArgs, id)
	return
}
