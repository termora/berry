package admin

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *Bot) addContributorCategory(ctx *bcr.Context) (err error) {
	name := ctx.Args[0]
	if len(name) > 200 {
		_, err = ctx.Replyc(bcr.ColourRed, "Name too long, maximum 200 characters (%v characters)", len(name))
		return
	}

	var role *discord.RoleID
	if len(ctx.Args) > 1 {
		r, err := ctx.ParseRole(ctx.Args[1])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find a role named `%v`", ctx.Args[1])
			return err
		}
		role = &r.ID
	}

	s := fmt.Sprintf("Are you sure you want to add a contributor category named \"%v\"", name)
	if role != nil {
		s += fmt.Sprintf(", with role %v", role.Mention())
	}
	s += "?"

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message: s,
	})
	if timeout {
		return ctx.SendX("Timed out.")
	}
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	cat, err := bot.DB.AddContributorCategory(name, role)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Reply("Added contributor category %v, with ID %v!", cat.Name, cat.ID)
	return
}

func (bot *Bot) listContributorCategories(ctx *bcr.Context) (err error) {
	cats, err := bot.DB.ContributorCategories()
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	if len(cats) == 0 {
		_, err = ctx.Reply("There are no contributor categories.")
		return
	}

	var s []string
	for _, cat := range cats {
		str := fmt.Sprintf("%v (%v)", cat.Name, cat.ID)
		if cat.RoleID != nil {
			str += fmt.Sprintf("\n%v", cat.RoleID.Mention())
		}
		str += "\n\n"

		s = append(s, str)
	}

	_, _, err = ctx.ButtonPages(
		bcr.StringPaginator("Contributor categories", db.EmbedColour, s, 5),
		5*time.Minute,
	)
	return
}

func (bot *Bot) addContributor(ctx *bcr.Context) (err error) {
	var (
		id   discord.UserID
		name string
	)

	m, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		u, err := ctx.ParseUser(ctx.Args[0])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "User not found.")
			return err
		}
		id = u.ID
		name = u.Username
	} else {
		id = m.User.ID
		name = m.Nick
		if m.Nick == "" {
			name = m.User.Username
		}
	}

	cat := bot.DB.ContributorCategory(ctx.Args[1])
	if cat == nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Category `%v` not found.", ctx.Args[1])
		return
	}

	err = bot.DB.AddContributor(cat.ID, id, name)
	if err != nil {
		// probably constraint error
		_, err = ctx.Replyc(bcr.ColourRed, "Error adding contributor: `%v`", err)
		return
	}

	_, err = ctx.Reply("Added %v as a %v!", id.Mention(), cat.Name)
	return
}

func (bot *Bot) overrideContributor(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find that user.")
		return err
	}

	name := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if name == ctx.RawArgs {
		name = strings.Join(ctx.Args[1:], " ")
	}

	var override *string
	if name != "clear" && name != "--clear" && name != "-clear" {
		override = &name
	}

	err = bot.DB.OverrideContributorName(u.ID, override)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Error updating override: %v", err)
		return
	}

	s := fmt.Sprintf("Cleared %v's name override!", u.Mention())
	if override != nil {
		s = fmt.Sprintf("Updated %v's name to \"%v\"!", u.Mention(), *override)
	}
	_, err = ctx.Reply(s)
	return
}

func (bot *Bot) allContributors(ctx *bcr.Context) (err error) {
	members, err := bot.Helper.Members(bot.Config.Bot.Support.GuildID, 0)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Error fetching members: %v", err)
		return
	}

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Embeds: []discord.Embed{{
			Color:       bcr.ColourRed,
			Description: fmt.Sprintf("This will check %v members' roles. Continue?", len(members)),
		}},
		YesPrompt: "Continue",
		NoPrompt:  "Cancel",
	})
	if timeout || !yes {
		_, err = ctx.Reply("Cancelled.")
		return
	}

	ctx.State.Typing(ctx.Channel.ID)

	// this is inefficient, but it should only run once anyway
	for _, m := range members {
		for _, r := range m.RoleIDs {
			cat := bot.DB.CategoryFromRole(r)
			if cat == nil {
				continue
			}

			name := m.User.Username
			if m.Nick != "" {
				name = m.Nick
			}

			err = bot.DB.AddContributor(cat.ID, m.User.ID, name)
			if err != nil {
				_, err = ctx.Replyc(bcr.ColourRed, "Error adding %v to category %v: %v", m.Mention(), cat.Name, err)
				return
			}
		}
	}

	_, err = ctx.Reply("Done!")
	return
}
