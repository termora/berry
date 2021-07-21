package search

import (
	"context"
	"strconv"
	"strings"
	"time"

	"flag"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) search(ctx *bcr.Context) (err error) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	var (
		showHidden bool
		noCW       bool
		cat        string
		ignore     string
		ignoreTags = []string{}
	)

	fs.BoolVar(&noCW, "no-cw", false, "")
	fs.BoolVar(&showHidden, "h", false, "")
	fs.StringVar(&cat, "c", "", "")
	fs.StringVar(&ignore, "i", "", "")

	fs.Parse(ctx.Args)
	ctx.Args = fs.Args()

	// we can't check for this normally because of the flags above
	// so we do it here, also lets us give a custom error message
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("No search term provided.")
		return err
	}

	// set tags to ignore
	if ignore != "" {
		ignoreTags = strings.Split(ignore, ",")
		for i := range ignoreTags {
			ignoreTags[i] = strings.ToLower(strings.TrimSpace(ignoreTags[i]))
		}
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
		terms, err = c.DB.Search(search, limit, ignoreTags)
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

		terms, err = c.DB.SearchCat(search, category, limit, showHidden, ignoreTags)
		if err != nil {
			return c.DB.InternalError(ctx, err)
		}
	}

	if len(terms) == 0 {
		_, err = ctx.Send("No results found.")
		return err
	}

	filter := []*db.Term{}
	if noCW {
		for _, t := range terms {
			if t.ContentWarnings == "" {
				filter = append(filter, t)
			}
		}
		terms = filter
	}

	// if there's only one term, just show that one
	if len(terms) == 1 {
		_, err = ctx.Send("", c.DB.TermEmbed(terms[0]))
		return err
	}

	// split the slice of terms into 5-long slices each
	var (
		termSlices [][]*db.Term
		embeds     []discord.Embed
	)

	for i := 0; i < len(terms); i += 5 {
		end := i + 5

		if end > len(terms) {
			end = len(terms)
		}

		termSlices = append(termSlices, terms[i:end])
	}

	// turn those slices into embeds
	for i, t := range termSlices {
		embeds = append(embeds, searchResultEmbed(search, i+1, len(termSlices), len(terms), t))
	}

	// actually send the search results
	msg, _, err := ctx.PagedEmbedTimeout(embeds, false, 15*time.Minute)
	if err != nil {
		c.Report(ctx, err)
		return err
	}

	// add the reactions, this is spun off into its own goroutine to immediately add the handler below
	go func() {
		for i, e := range emoji {
			if i >= len(terms) {
				return
			}

			err = ctx.State.React(ctx.Channel.ID, msg.ID, discord.APIEmoji(e))
			// if the error was non-nil we can assume the message was deleted, so return
			if err != nil {
				break
			}
		}
	}()

	// time out the request below after 15 minutes
	// deferring the cancel func is just good practice
	con, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// wait for either a message or a reaction
	// store the number in a variable so we don't have to parse it all over again
	var n int
	v := ctx.State.WaitFor(con, func(v interface{}) bool {
		if ev, ok := v.(*gateway.MessageCreateEvent); ok {
			// if the channel/author aren't correct, return
			if ev.Author.ID != ctx.Author.ID || ev.ChannelID != ctx.Message.ChannelID {
				return false
			}

			// parse the number
			n, err = strconv.Atoi(ev.Content)
			if err != nil {
				return false
			}
		} else if ev, ok := v.(*gateway.MessageReactionAddEvent); ok {
			// else, check for a message reaction
			if ev.UserID != ctx.Author.ID ||
				ev.ChannelID != ctx.Message.ChannelID ||
				ev.MessageID != msg.ID {
				return false
			}

			// get the emoji number
			var isNum bool
			for i, e := range emoji {
				if ev.Emoji.Name == e {
					n = i + 1
					isNum = true
					break
				}
			}

			// if it wasn't a number emoji, return
			if !isNum {
				return false
			}
		} else {
			return false
		}

		// get the page number
		// this conversion *shouldn't* fail, but if it does and we don't check for that, the function will panic
		page, ok := ctx.AdditionalParams["page"].(int)
		if !ok {
			return false
		}

		// this should never happen but check just in case
		if len(termSlices) < page {
			return false
		}

		// if the reaction/number is out of bounds, return
		if n > len(termSlices[page]) {
			return false
		}

		// everything's fine so we can accept this event!
		return true
	})

	// if it timed out, return
	// and try to clean up reactions too
	if v == nil {
		return
	}

	page, ok := ctx.AdditionalParams["page"].(int)
	if !ok {
		return
	}

	// delete the original message, then send the definition
	ctx.State.DeleteMessage(ctx.Channel.ID, msg.ID)
	_, err = ctx.Send("", c.DB.TermEmbed(termSlices[page][n-1]))
	return
}
