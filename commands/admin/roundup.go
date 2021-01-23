package admin

import (
	"fmt"
	"strings"
	"time"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
)

func (c *commands) changelog(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(2); err != nil {
		_, err = ctx.Send("You are missing the required arguments `channel` and/or `since`.", nil)
		return err
	}

	// parse channel
	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("That channel could not be found.", nil)
		return
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("That channel is not in this server.", nil)
		return
	}

	// parse date
	date, err := time.Parse("2006-01-02", ctx.Args[1])
	if err != nil {
		_, err = ctx.Send("Please input the date as `yyyy-mm-dd`.", nil)
		return err
	}

	t, err := c.db.TermsSince(date)
	if err != nil {
		return c.db.InternalError(ctx, err)
	}

	if len(t) == 0 {
		_, err = ctx.Sendf("No terms were added since the specified date (%v).", date.Format("January 2, 2006"))
		return err
	}

	terms := make([]string, 0)
	for _, term := range t {
		terms = append(terms, term.Name)
	}

	// check perms in the channel
	perms, err := ctx.Session.Permissions(ch.ID, ctx.Author.ID)
	if err != nil {
		c.sugar.Errorf("Error getting perms for %v in %v: %v", ctx.Author.ID, ch.ID, err)
		_, err = ctx.Sendf(
			"❌ An error occurred while trying to get permissions.\nIf this issue persists, please contact the bot developer.",
		)
		return
	}
	if perms&discord.PermissionSendMessages != discord.PermissionSendMessages || perms&discord.PermissionViewChannel != discord.PermissionViewChannel {
		_, err = ctx.Sendf(
			"❌ Error: this command requires the `%v` permissions in the channel you're posting to.",
			strings.Join(bcr.PermStrings(discord.PermissionSendMessages|discord.PermissionViewChannel), ", "),
		)
		return
	}

	msgs := make([]string, 0)
	s := fmt.Sprintf(
		"Since %v, **%v** new terms have been added, for a total of **%v** terms!\n\n**New terms**\nThe following terms have been added: %v",
		date.Format("January 02"), len(t),
		c.db.TermCount(), strings.Join(terms, ", "),
	)
	if len(s) >= 2000 {
		s = fmt.Sprintf("Since %v, **%v** new terms have been added, for a total of **%v** terms!", date.Format("January 02"), len(t), c.db.TermCount())

		buf := "**New terms**\nThe following terms have been added:\n"
		for _, t := range terms {
			if len(buf) >= 1900 {
				msgs = append(msgs, buf)
				buf = ""
			}
			buf += t + ", "
		}
		msgs = append(msgs, buf)
	}

	_, err = ctx.Session.SendMessage(ch.ID, c.config.Bot.TermChangelogPing, &discord.Embed{
		Title:       "Term changelog",
		Description: s,

		Color: db.EmbedColour,
	})
	if err != nil {
		return err
	}

	if len(msgs) > 0 {
		for _, m := range msgs {
			time.Sleep(500 * time.Millisecond)
			_, err = ctx.Session.SendMessage(ch.ID, "", &discord.Embed{
				Description: m,
				Color:       db.EmbedColour,
			})
			if err != nil {
				return err
			}
		}
	}
	return
}
