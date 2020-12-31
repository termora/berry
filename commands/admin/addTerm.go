package admin

import (
	"strings"

	"github.com/Starshine113/crouter"
	"github.com/Starshine113/berry/db"
	"github.com/bwmarrin/discordgo"
)

func (c *commands) addTerm(ctx *crouter.Ctx) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		return ctx.CommandError(err)
	}

	term := &db.Term{Name: ctx.RawArgs}
	ctx.AdditionalParams["term"] = term

	m, err := ctx.Sendf("Creating a term with the name `%v`. To cancel at any time, send `cancel`.\nPlease type the name of the category this term belongs to:", ctx.RawArgs)
	if err != nil {
		return err
	}

	ctx.AddMessageHandler(m.ID, func(ctx *crouter.Ctx, m *discordgo.MessageCreate) {
		if m.Content == "cancel" {
			ctx.Send("Term creation cancelled.")
			return
		}
		cat, err := c.db.CategoryID(m.Content)
		if err != nil {
			ctx.CommandError(err)
			return
		}
		if cat == 0 {
			return
		}

		t := ctx.AdditionalParams["term"].(*db.Term)
		t.Category = cat
		ctx.AdditionalParams["term"] = t
		msg, err := ctx.Sendf("Category set to `%v` (ID %v). Please type the description:", m.Content, cat)
		if err != nil {
			return
		}

		ctx.AddMessageHandler(msg.ID, func(ctx *crouter.Ctx, m *discordgo.MessageCreate) {
			if m.Content == "cancel" {
				ctx.Send("Term creation cancelled.")
				return
			}
			t := ctx.AdditionalParams["term"].(*db.Term)
			t.Description = m.Content
			ctx.AdditionalParams["term"] = t
			msg, err := ctx.Send("Description set. Please type the source:")
			if err != nil {
				return
			}

			ctx.AddMessageHandler(msg.ID, func(ctx *crouter.Ctx, m *discordgo.MessageCreate) {
				if m.Content == "cancel" {
					ctx.Send("Term creation cancelled.")
					return
				}
				t := ctx.AdditionalParams["term"].(*db.Term)
				t.Source = m.Content
				ctx.AdditionalParams["term"] = t
				msg, err := ctx.Send("Source set. Please type a *newline separated* list of aliases/synonyms, or \"none\" to set no aliases:")
				if err != nil {
					return
				}

				ctx.AddMessageHandler(msg.ID, func(ctx *crouter.Ctx, m *discordgo.MessageCreate) {
					if m.Content == "cancel" {
						ctx.Send("Term creation cancelled.")
						return
					}
					t := ctx.AdditionalParams["term"].(*db.Term)
					t.Aliases = strings.Split(m.Content, "\n")
					if m.Content == "none" {
						t.Aliases = []string{}
					}

					msg, err := ctx.Send(&discordgo.MessageSend{
						Content: "Term finished. React with ✅ to finish adding it, or with ❌ to cancel. Preview:",
						Embed:   t.TermEmbed(),
					})
					if err != nil {
						return
					}

					ctx.AddYesNoHandler(msg.ID, func(ctx *crouter.Ctx) {
						t, err := c.db.AddTerm(t)
						if err != nil {
							ctx.CommandError(err)
							return
						}
						ctx.Sendf("Added term with ID %v.", t.ID)
					}, func(ctx *crouter.Ctx) {
						ctx.Send("Cancelled.")
					})
				})
			})
		})
	})

	return
}
