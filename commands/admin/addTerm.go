package admin

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/admin/auditlog"
	"github.com/termora/berry/db"
)

func (bot *Bot) addTerm(ctx *bcr.Context) (err error) {
	t := &db.Term{}

	names := strings.Split(ctx.RawArgs, "\n")
	t.Name = names[0]
	if len(names) > 1 {
		t.Aliases = names[1:]
	}

	embed := discord.Embed{
		Title: t.Name,
		Color: db.EmbedColour,
	}
	if len(t.Aliases) > 0 {
		embed.Fields = append(embed.Fields, discord.EmbedField{
			Name:  "Synonyms",
			Value: strings.Join(t.Aliases, ", "),
		})
	}

	info, err := ctx.Send(
		"Adding a new term; to cancel at any time, type `cancel`.\nPlease send the description.",
		embed,
	)
	if err != nil {
		return
	}

	m, timeout := ctx.WaitForMessage(ctx.Channel.ID, ctx.Author.ID, 15*time.Minute, nil)
	if timeout {
		_, err = ctx.Send(":x: Timed out.")
		return
	}
	if strings.EqualFold(m.Content, "cancel") {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	t.Description = m.Content
	embed.Description = m.Content

	_, err = ctx.Edit(info, "Adding a new term; to cancel at any time, type `cancel`.\nPlease send the source.", true, embed)
	if err != nil {
		return err
	}
	ctx.State.DeleteMessage(m.ChannelID, m.ID, "")

	m, timeout = ctx.WaitForMessage(ctx.Channel.ID, ctx.Author.ID, 15*time.Minute, nil)
	if timeout {
		_, err = ctx.Send(":x: Timed out.")
		return
	}
	if strings.EqualFold(m.Content, "cancel") {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	t.Source = m.Content

	_, err = ctx.Edit(info, "Adding a new term; to cancel at any time, type `cancel`.\nPlease send a list of tags, separated by newlines.\n__Note that the first tag needs to be a valid category name__.", true, bot.DB.TermEmbed(t))
	if err != nil {
		return err
	}
	ctx.State.DeleteMessage(m.ChannelID, m.ID, "")

	m, timeout = ctx.WaitForMessage(ctx.Channel.ID, ctx.Author.ID, 15*time.Minute, nil)
	if timeout {
		_, err = ctx.Send(":x: Timed out.")
		return
	}
	if strings.EqualFold(m.Content, "cancel") {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	tags := strings.Split(m.Content, "\n")
	category, err := bot.DB.CategoryID(tags[0])
	if err != nil {
		_, err = ctx.Sendf(":x: Couldn't find a category with that name (``%v``).", tags[0])
	}

	t.Category = category
	t.CategoryName = tags[0]
	t.DisplayTags = tags
	for _, tag := range tags {
		t.Tags = append(t.Tags, strings.ToLower(strings.TrimSpace(tag)))

		con, cancel := bot.DB.Context()
		defer cancel()

		_, err = bot.DB.Exec(con, `insert into public.tags (normalized, display) values ($1, $2)
		on conflict (normalized) do update set display = $2`, strings.ToLower(strings.TrimSpace(tag)), tag)
		if err != nil {
			bot.Log.Errorf("Error adding tag: %v", err)
		}
	}

	_, err = ctx.Edit(info, "Are you sure you want to add this term?", true, bot.DB.TermEmbed(t))
	if err != nil {
		return err
	}
	ctx.State.DeleteMessage(m.ChannelID, m.ID, "")

	yes, timeout := ctx.YesNoHandler(*info, ctx.Author.ID)
	if timeout {
		_, err = ctx.Send(":x: Operation timed out.")
		return
	}
	if !yes {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	t, err = bot.DB.AddTerm(t)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = bot.AuditLog.SendLog(t.ID, auditlog.TermEntry, auditlog.CreateAction, nil, t, ctx.Author.ID, nil)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Added term with ID %v.", t.ID)
	return
}
