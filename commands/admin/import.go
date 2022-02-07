package admin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/commands/admin/auditlog"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
)

var msgRegex = regexp.MustCompile(`(?i)\*?\*?Name:?\*?\*?:?(.*)\n\*?\*?Category:?\*?\*?:?(.*)\n\*?\*?Description:?\*?\*?:?([\s\S]*)\n\*?\*?Coined by:?\*?\*?:?(.*)`)

var tagsRegex = regexp.MustCompile(`(?i)\*?\*?Tags:\*?\*? (.*)`)

func (bot *Bot) importFromMessage(ctx *bcr.Context) (err error) {
	flag, _ := ctx.Flags.GetString("category")
	rawSource, _ := ctx.Flags.GetBool("raw-source")

	msg, err := ctx.ParseMessage(ctx.Args[0])
	if err != nil {
		_, err = ctx.Sendf("Message not found. Are you sure I have access to that channel?")
		return
	}

	t := &db.Term{}

	// embeds are easy, just parse all of the fields
	if len(msg.Embeds) > 0 {
		// only use the first embed
		for _, f := range msg.Embeds[0].Fields {
			if f.Name == "Term title" {
				t.Name = f.Value
				continue
			}
			switch f.Name {
			case "Term title":
				t.Name = f.Value
			case "Aliases/other names (optional, comma-separated)":
				aliases := strings.Split(f.Value, ",")
				for i := range aliases {
					aliases[i] = strings.TrimSpace(aliases[i])
				}
				t.Aliases = aliases
			case "Tags":
				tags := strings.Split(f.Value, ",")
				for i := range tags {
					tags[i] = strings.TrimSpace(tags[i])
				}
				t.Tags = tags
			case "Description":
				t.Description = f.Value
			case "Source":
				t.Source = f.Value
			case "What category does your term fall under? Pick the most relevant one.":
				if cat, err := bot.DB.CategoryID(f.Value); err == nil {
					t.Category = cat
					t.CategoryName = f.Value
				}
			}
		}

		// we're done parsing the term
		goto done
	}

	// otherwise we'll have to parse the content
	if !msgRegex.MatchString(msg.Content) {
		ctx.SendfX("⚠️ The message you gave didn't match the expected input. You might have to add it manually with ``%vadmin addterm``.", bot.Config.Bot.Prefixes[0])

		// the message didn't match, so don't bother parsing everything
		goto done
	}

	{
		groups := msgRegex.FindStringSubmatch(msg.Content)

		// names
		names := strings.Split(groups[1], ",")
		for i := range names {
			names[i] = strings.TrimSpace(names[i])
		}
		t.Name = names[0]
		if len(names) > 1 {
			t.Aliases = names[1:]
		}

		// category
		cat, err := bot.DB.CategoryID(strings.TrimSpace(groups[2]))
		if err == nil {
			t.Category = cat
			t.CategoryName = groups[2]
		}

		t.Description = strings.TrimSpace(groups[3])
		t.Source = strings.TrimSpace(groups[4])

		if g := tagsRegex.FindStringSubmatch(msg.Content); len(g) > 1 {
			tags := strings.Split(g[1], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			t.Tags = tags
		}
	}

done:
	// validate the term object
	if t.Name == "" || t.Source == "" || t.Description == "" {
		return ctx.SendfX("One or more required fields (name, source, description) was empty!\n(Debug: name length => %v, source length => %v, description length => %v)", len(t.Name), len(t.Source), len(t.Description))
	}
	if t.Aliases == nil {
		t.Aliases = []string{}
	}
	if !rawSource && !bcr.HasAnyPrefix(t.Source, "Coined by", "Unknown", "unknown", "Already") {
		t.Source = fmt.Sprintf("Coined by %v", t.Source)
	}

	if t.Category == 0 {
		if flag == "" {
			_, err = ctx.Send("No category specified, and the submission didn't specify a category.")
			return
		}

		cat, err := bot.DB.CategoryID(flag)
		if err != nil {
			_, err = ctx.Sendf("That category (``%v``) could not be found.", bcr.EscapeBackticks(flag))
		}
		t.Category = cat
	}

	// add the category to the tags, if it's not already in there
	cat := bot.DB.CategoryFromID(t.Category)
	t.CategoryName = cat.Name

	catInTags := false
	for _, tag := range t.Tags {
		if tag == cat.Name {
			catInTags = true
			break
		}
	}
	if !catInTags {
		t.Tags = append(t.Tags, cat.Name)
	}

	// these aren't used when inserting the term, just for TermEmbed below
	t.DisplayTags = t.Tags

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message:   "Do you want to add this term?",
		Embeds:    []discord.Embed{bot.DB.TermEmbed(t)},
		YesPrompt: "Add term",
		YesStyle:  discord.SuccessButtonStyle(),
	})
	if timeout {
		_, err = ctx.Send(":x: Operation timed out.")
		return
	}
	if !yes {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	for i := range t.Tags {
		con, cancel := bot.DB.Context()
		defer cancel()

		_, err = bot.DB.Exec(con, `insert into public.tags (normalized, display) values ($1, $2)
		on conflict (normalized) do update set display = $2`, strings.ToLower(t.Tags[i]), t.Tags[i])
		if err != nil {
			log.Errorf("Error adding tag: %v", err)
		}

		t.DisplayTags = append(t.DisplayTags, t.Tags[i])
		t.Tags[i] = strings.ToLower(t.Tags[i])
	}

	t, err = bot.DB.AddTerm(t)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	_, err = ctx.Sendf("Added term with ID %v.", t.ID)
	if err != nil {
		return err
	}

	// if we don't have perms return
	if p, _ := ctx.State.Permissions(msg.ChannelID, ctx.Bot.ID); !p.Has(discord.PermissionAddReactions | discord.PermissionReadMessageHistory) {
		return
	}

	// react with a checkmark to the original message
	ctx.State.React(msg.ChannelID, msg.ID, "yes:822929172669136966")

	// if logging terms is enabled, log this
	_, err = bot.AuditLog.SendLog(t.ID, auditlog.TermEntry, auditlog.CreateAction, nil, t, ctx.Author.ID, nil)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}
	return err
}
