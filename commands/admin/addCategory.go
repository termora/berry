package admin

import (
	"context"

	"github.com/Starshine113/crouter"
)

func (c *commands) addCategory(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	var e bool
	err = c.db.Pool.QueryRow(context.Background(), "select exists (select from categories where lower(name) = lower($1))", ctx.RawArgs).Scan(&e)
	if err != nil {
		return ctx.CommandError(err)
	}
	if e {
		_, err = ctx.Sendf("%v A category with that name already exists.", crouter.ErrorEmoji)
		return err
	}

	var id int
	err = c.db.Pool.QueryRow(context.Background(), "insert into public.categories (name) values ($1) returning id", ctx.RawArgs).Scan(&id)
	if err != nil {
		return ctx.CommandError(err)
	}
	_, err = ctx.Sendf("Added category `%v` with ID %v.", ctx.RawArgs, id)
	return
}
