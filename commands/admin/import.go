package admin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

var msgRegex = regexp.MustCompile(`(?i)\*?\*?Name:?\*?\*?:?(.*)\n\*?\*?Category:?\*?\*?:?(.*)\n\*?\*?Description:?\*?\*?:?([\s\S]*)\n\*?\*?Coined by:?\*?\*?:?(.*)`)

var tagsRegex = regexp.MustCompile(`(?i)\*?\*?Tags:\*?\*? (.*)`)

func (c *Admin) importFromMessage(ctx *bcr.Context) (err error) {
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
				if cat, err := c.DB.CategoryID(f.Value); err == nil {
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
		cat, err := c.DB.CategoryID(strings.TrimSpace(groups[2]))
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
		_, err = ctx.Send("One or more required fields (name, source, description) was empty!")
		return
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

		cat, err := c.DB.CategoryID(flag)
		if err != nil {
			_, err = ctx.Sendf("That category (``%v``) could not be found.", bcr.EscapeBackticks(flag))
		}
		t.Category = cat
	}

	// add the category to the tags, if it's not already in there
	cat := c.DB.CategoryFromID(t.Category)
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

	termMsg, err := ctx.Send("Do you want to add this term?", t.TermEmbed(c.Config.TermBaseURL()))
	if err != nil {
		return err
	}

	yes, timeout := ctx.YesNoHandler(*termMsg, ctx.Author.ID)
	if timeout {
		_, err = ctx.Send(":x: Operation timed out.")
		return
	}
	if !yes {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	for i := range t.Tags {
		con, cancel := c.DB.Context()
		defer cancel()

		_, err = c.DB.Pool.Exec(con, `insert into public.tags (normalized, display) values ($1, $2)
		on conflict (normalized) do update set display = $2`, strings.ToLower(t.Tags[i]), t.Tags[i])
		if err != nil {
			c.Sugar.Errorf("Error adding tag: %v", err)
		}

		t.DisplayTags = append(t.DisplayTags, t.Tags[i])
		t.Tags[i] = strings.ToLower(t.Tags[i])
	}

	t, err = c.DB.AddTerm(t)
	if err != nil {
		return c.DB.InternalError(ctx, err)
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
	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Content: "​",

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Description: fmt.Sprintf("Term imported from\n%v", ctx.Args[0]),
					Color:       db.EmbedColour,
					Timestamp:   discord.NowTimestamp(),
				},
				e,
			},
		})
	}
	return err
}
