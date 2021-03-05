package search

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	flag "github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) search(ctx *bcr.Context) (err error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	var (
		showHidden bool
		cat        string
	)

	fs.BoolVarP(&showHidden, "show-hidden", "h", false, "")
	fs.StringVarP(&cat, "category", "c", "", "")

	err = fs.Parse(ctx.Args)
	if err != nil {
		_, err = ctx.Sendf("You didn't give a valid input for this command. Usage:\n> ``%v%v %v``", ctx.Router.Prefixes[0], ctx.Command, ctx.Cmd.Usage)
		return
	}
	ctx.Args = fs.Args()

	// we can't check for this normally because of the flags above
	// so we do it here, also lets us give a custom error message
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("No search term provided.", nil)
		return err
	}

	search := strings.Join(ctx.Args, " ")

	limit := 0
	// if the query starts with !, only show the first result
	if strings.HasPrefix(search, "!") {
		limit = 1
		search = strings.TrimPrefix(search, "!")
	}

	var terms []*db.Term
	if cat == "" {
		// no category given, so just search *all* terms
		terms, err = c.DB.Search(search, limit)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
	} else {
		// category given, so search in category

		// get the category ID
		category, err := c.DB.CategoryID(cat)
		if err != nil {
			_, err = ctx.Sendf("The category you specified (``%v``) was not found.", bcr.EscapeBackticks(cat))
			return err
		}

		terms, err = c.DB.SearchCat(search, category, limit, showHidden)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
	}

	if len(terms) == 0 {
		_, err = ctx.Send("No results found.", nil)
		return err
	}
	// if there's only one term, just show that one
	if len(terms) == 1 {
		_, err = ctx.Send("", terms[0].TermEmbed(c.Config.TermBaseURL()))
		return err
	}

	termSlices := make([][]*db.Term, 0)

	for i := 0; i < len(terms); i += 5 {
		end := i + 5

		if end > len(terms) {
			end = len(terms)
		}

		termSlices = append(termSlices, terms[i:end])
	}

	embeds := make([]discord.Embed, 0)

	for i, t := range termSlices {
		embeds = append(embeds, searchResultEmbed(search, i+1, len(termSlices), len(terms), t))
	}

	msg, err := ctx.PagedEmbed(embeds, false)
	if err != nil {
		c.Report(ctx, err)
		return err
	}

	ctx.AdditionalParams["termSlices"] = termSlices

	for i, e := range emoji {
		if i >= len(terms) {
			return
		}

		emoji := e
		if err := ctx.Session.React(ctx.Channel.ID, msg.ID, discord.APIEmoji(emoji)); err != nil {
			c.Sugar.Error("Error adding reaction:", err)
			return err
		}

		index := i
		ctx.AddReactionHandler(msg.ID, ctx.Author.ID, e, false, false, func(ctx *bcr.Context) {
			page, ok := ctx.AdditionalParams["page"].(int)
			if ok == false {
				return
			}
			termSlices, ok := ctx.AdditionalParams["termSlices"].([][]*db.Term)
			if ok == false {
				return
			}
			if len(termSlices) < page {
				ctx.Session.DeleteUserReaction(ctx.Channel.ID, msg.ID, ctx.Author.ID, discord.APIEmoji(emoji))
				return
			}

			termSlice := termSlices[page]
			if index >= len(termSlice) {
				ctx.Session.DeleteUserReaction(ctx.Channel.ID, msg.ID, ctx.Author.ID, discord.APIEmoji(emoji))
				return
			}

			err := ctx.Session.DeleteMessage(ctx.Channel.ID, msg.ID)
			if err != nil {
				c.Sugar.Error("Error deleting message:", err)
			}
			_, err = ctx.Send("", termSlice[index].TermEmbed(c.Config.TermBaseURL()))
			if err != nil {
				c.Sugar.Error("Error sending message:", err)
			}
		})
	}

	return
}
