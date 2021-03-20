package admin

import (
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

var msgRegex = regexp.MustCompile(`\*\*Name:\*\* (.*)\n\*\*Category:\*\* (.*)\n\*\*Description:\*\* ([\s\S]*)\n\*\*Coined by:\*\* (.*)`)

func (c *Admin) importFromMessage(ctx *bcr.Context) (err error) {
	var flag string

	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVarP(&flag, "category", "c", "", "Category")
	fs.Parse(ctx.Args)
	ctx.Args = fs.Args()

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
		cat, err := c.DB.CategoryID(groups[2])
		if err == nil {
			t.Category = cat
			t.CategoryName = groups[2]
		}

		t.Description = groups[3]
		t.Source = groups[4]
	}

done:
	// validate the term object
	if t.Name == "" || t.Source == "" || t.Description == "" {
		_, err = ctx.Send("One or more required fields (name, source, description) was empty!", nil)
		return
	}
	if t.Aliases == nil {
		t.Aliases = []string{}
	}

	if t.Category == 0 {
		if flag == "" {
			_, err = ctx.Send("No category specified, and the submission didn't specify a category.", nil)
			return
		}

		cat, err := c.DB.CategoryID(flag)
		if err != nil {
			_, err = ctx.Sendf("That category (``%v``) could not be found.", bcr.EscapeBackticks(flag))
		}
		t.Category = cat
		t.CategoryName = flag
	}

	termMsg, err := ctx.Send("Do you want to add this term?", t.TermEmbed(c.Config.TermBaseURL()))
	if err != nil {
		return err
	}

	yes, timeout := ctx.YesNoHandler(*termMsg, ctx.Author.ID)
	if timeout {
		_, err = ctx.Send(":x: Operation timed out.", nil)
		return
	}
	if !yes {
		_, err = ctx.Send(":x: Cancelled.", nil)
		return
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
	if p, _ := ctx.Session.Permissions(msg.ChannelID, ctx.Bot.ID); !p.Has(discord.PermissionAddReactions | discord.PermissionReadMessageHistory) {
		return
	}

	// react with a checkmark to the original message
	err = ctx.Session.React(msg.ChannelID, msg.ID, "yes:822929172669136966")
	return err
}
