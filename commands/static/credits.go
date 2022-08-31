package static

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *Bot) credits(ctx bcr.Contexter) (err error) {
	bot.memberMu.RLock()
	defer bot.memberMu.RUnlock()

	// return if there's no credit fields
	if len(bot.Config.Bot.CreditFields) == 0 &&
		(len(bot.Config.Bot.ContributorRoles) == 0 ||
			len(bot.SupportServerMembers) == 0) {
		return nil
	}

	embeds := []discord.Embed{{
		Color:       db.EmbedColour,
		Title:       "Credits",
		Description: fmt.Sprintf("These are the people who helped create %v!", bot.Router.Bot.Username),
		Fields:      bot.Config.Bot.CreditFields,
	}}

	e := discord.Embed{
		Color:       db.EmbedColour,
		Title:       "Contributors",
		Description: fmt.Sprintf("These are the people who have contributed to %v in some capacity!", bot.Router.Bot.Username),
	}

	cats, err := bot.DB.ContributorCategories()
	if err != nil {
		_, _, err = ctx.ButtonPages(embeds, 15*time.Minute)
		return err
	}

	for _, cat := range cats {
		contributors, err := bot.DB.Contributors(cat.ID)
		if err != nil {
			_, _, err = ctx.ButtonPages(embeds, 15*time.Minute)
			return err
		}

		if len(contributors) == 0 {
			continue
		}

		var (
			slice []string
			s     string
		)
		for _, m := range contributors {
			name := m.Name
			if m.Override != nil {
				name = *m.Override
			}
			slice = append(slice, name)
		}
		for i, m := range slice {
			if len(s) > 900 {
				s += fmt.Sprintf("\n...and %v others!", len(slice)-i)
				break
			}
			if i != 0 {
				s += ", "
			}
			s += m
		}

		name := cat.Name
		if len(slice) != 1 {
			name += "s"
		}
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:   name,
			Value:  s,
			Inline: false,
		})
	}
	if len(e.Fields) > 0 {
		embeds = append(embeds, e)
		embeds[0].Description += "\nReact with ➡️ to show everyone who has contributed!"
	}

	_, _, err = ctx.ButtonPages(embeds, 15*time.Minute)
	return err
}
