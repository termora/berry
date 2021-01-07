package search

import (
	"time"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/Starshine113/berry/misc"
	"github.com/Starshine113/berry/structs"
	"github.com/diamondburned/arikawa/v2/discord"
	"go.uber.org/zap"
)

type commands struct {
	Db    *db.Db
	Sugar *zap.SugaredLogger
	conf  *structs.BotConfig
}

// Init ...
func Init(db *db.Db, conf *structs.BotConfig, s *zap.SugaredLogger, r *bcr.Router) {
	c := commands{Db: db, conf: conf, Sugar: s}

	r.AddCommand(&bcr.Command{
		Name:    "search",
		Aliases: []string{"s"},

		Description: "Search for a term",
		Usage:       "<search term>",

		Cooldown:      3 * time.Second,
		Blacklistable: true,

		Command: c.search,
	})

	r.AddCommand(&bcr.Command{
		Name:    "random",
		Aliases: []string{"r"},

		Description: "Show a random term",

		Cooldown:      3 * time.Second,
		Blacklistable: true,

		Command: c.random,
	})

	r.AddCommand(&bcr.Command{
		Name:    "explain",
		Aliases: []string{"e", "ex"},

		Description: "",
		Usage:       "[explanation]",

		Cooldown:      1 * time.Second,
		Blacklistable: false,

		Command: c.explanation,
	})

	r.AddCommand(&bcr.Command{
		Name:        "list",
		Description: "List all terms",

		Cooldown:      3 * time.Second,
		Blacklistable: true,
		Command:       c.list,
	})

	r.AddCommand(&bcr.Command{
		Name:        "post",
		Description: "Post a single term",
		Usage:       "<term ID> [channel]",

		Cooldown:      3 * time.Second,
		Blacklistable: true,
		Command:       c.term,
	})
}

func (c *commands) search(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("No search term provided.", nil)
		return err
	}

	terms, err := c.Db.Search(ctx.RawArgs, 0)
	if err != nil {
		_, err = ctx.Send(misc.InternalError, nil)
		return err
	}

	if len(terms) == 0 {
		_, err = ctx.Send("No results found.", nil)
		return err
	}
	if len(terms) == 1 {
		_, err = ctx.Send("", terms[0].TermEmbed(c.conf.Bot.TermBaseURL))
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
		embeds = append(embeds, searchResultEmbed(ctx.RawArgs, i+1, len(termSlices), t))
	}

	msg, err := ctx.PagedEmbed(embeds)
	if err != nil {
		return err
	}

	ctx.AdditionalParams["termSlices"] = termSlices

	for i, e := range emoji {
		emoji := e
		if err = ctx.Session.React(ctx.Channel.ID, msg.ID, discord.APIEmoji(emoji)); err != nil {
			return
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

			ctx.Session.DeleteMessage(ctx.Channel.ID, msg.ID)
			ctx.Send("", termSlice[index].TermEmbed(c.conf.Bot.TermBaseURL))
		})
	}

	return
}
