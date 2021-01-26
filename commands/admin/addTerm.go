package admin

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
)

func (c *Admin) addTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("Please provide a term name.", nil)
		return err
	}

	term := &db.Term{Name: ctx.RawArgs}
	ctx.AdditionalParams["term"] = term

	_, err = ctx.Sendf("Creating a term with the name `%v`. To cancel at any time, send `cancel`.\nPlease type the name of the category this term belongs to:", ctx.RawArgs)
	if err != nil {
		return err
	}

	ctx.AddMessageHandler(ctx.Channel.ID, ctx.Author.ID, func(ctx *bcr.Context, m discord.Message) {
		if m.Content == "cancel" {
			ctx.Send("Term creation cancelled.", nil)
			return
		}
		cat, err := c.db.CategoryID(m.Content)
		if err != nil {
			_, err = ctx.Send("Could not find that category, cancelled.", nil)
			return
		}
		if cat == 0 {
			return
		}

		t := ctx.AdditionalParams["term"].(*db.Term)
		t.Category = cat
		ctx.AdditionalParams["term"] = t
		_, err = ctx.Sendf("Category set to `%v` (ID %v). Please type the description:", m.Content, cat)
		if err != nil {
			return
		}

		ctx.AddMessageHandler(ctx.Channel.ID, ctx.Author.ID, func(ctx *bcr.Context, m discord.Message) {
			if m.Content == "cancel" {
				ctx.Send("Term creation cancelled.", nil)
				return
			}
			t := ctx.AdditionalParams["term"].(*db.Term)
			t.Description = m.Content
			if len(t.Description) > 1800 {
				_, err = ctx.Send("Description too long (maximum 1800 characters).", nil)
				return
			}
			ctx.AdditionalParams["term"] = t
			_, err := ctx.Send("Description set. Please type the source:", nil)
			if err != nil {
				return
			}

			ctx.AddMessageHandler(ctx.Channel.ID, ctx.Author.ID, func(ctx *bcr.Context, m discord.Message) {
				if m.Content == "cancel" {
					ctx.Send("Term creation cancelled.", nil)
					return
				}
				t := ctx.AdditionalParams["term"].(*db.Term)
				t.Source = m.Content
				ctx.AdditionalParams["term"] = t
				_, err := ctx.Send("Source set. Please type a *newline separated* list of aliases/synonyms, or \"none\" to set no aliases:", nil)
				if err != nil {
					return
				}

				ctx.AddMessageHandler(ctx.Channel.ID, ctx.Author.ID, func(ctx *bcr.Context, m discord.Message) {
					if m.Content == "cancel" {
						ctx.Send("Term creation cancelled.", nil)
						return
					}
					t := ctx.AdditionalParams["term"].(*db.Term)
					t.Aliases = strings.Split(m.Content, "\n")
					if m.Content == "none" {
						t.Aliases = []string{}
					}

					msg, err := ctx.Send("Term finished. React with ✅ to finish adding it, or with ❌ to cancel. Preview:", t.TermEmbed(""))
					if err != nil {
						return
					}

					ctx.AddYesNoHandler(*msg, ctx.Author.ID, func(ctx *bcr.Context) {
						t, err := c.db.AddTerm(t)
						if err != nil {
							c.db.InternalError(ctx, err)
							return
						}
						ctx.Sendf("Added term with ID %v.", t.ID)
					}, func(ctx *bcr.Context) {
						ctx.Send("Cancelled.", nil)
					})
				})
			})
		})
	})

	return
}
