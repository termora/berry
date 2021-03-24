package admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *Admin) editTerm(ctx *bcr.Context) (err error) {
	if err = ctx.CheckMinArgs(3); err != nil {
		_, err = ctx.Send("Not enough arguments. Valid subcommands are: `title`, `desc`, `source`, `aliases`, `image`, `tags`.", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Args[1])
	if err != nil {
		_, err = ctx.Sendf("Could not parse ID:\n```%v```", err)
		return
	}
	t, err := c.DB.GetTerm(id)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("No term with that ID found.", nil)
			return
		}

		return c.DB.InternalError(ctx, err)
	}

	// these should probably be actual subcommands but then we'd have to duplicate the code above 6 times
	switch ctx.Args[0] {
	case "name", "title":
		return c.editTermTitle(ctx, t)
	case "desc", "description":
		return c.editTermDesc(ctx, t)
	case "source":
		return c.editTermSource(ctx, t)
	case "image":
		return c.editTermImage(ctx, t)
	case "tags":
		return c.editTermTags(ctx, t)
	case "aliases":
		return c.editTermAliases(ctx, t)
	}

	_, err = ctx.Send("Invalid subcommand supplied.\nValid subcommands are: `title`, `desc`, `source`, `aliases`, `image`.", nil)
	return
}

func (c *Admin) editTermTitle(ctx *bcr.Context, t *db.Term) (err error) {
	title := strings.Join(ctx.Args[2:], " ")
	if len(title) > 200 {
		_, err = ctx.Sendf("Title too long (%v > 200).", len(title))
		return
	}

	err = c.DB.UpdateTitle(t.ID, title)
	if err != nil {
		_, err = ctx.Sendf("Error updating title: ```%v```", err)
		return
	}

	_, err = ctx.Send("Title updated!", nil)
	if err != nil {
		c.Report(ctx, err)
	}

	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		e.Author = &discord.EmbedAuthor{
			Name: "Previous version",
		}

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title:       "Title updated",
					Description: title,
					Color:       db.EmbedColour,
					Timestamp:   discord.NowTimestamp(),
				},
				*e,
			},
		})
	}
	return
}

func (c *Admin) editTermDesc(ctx *bcr.Context, t *db.Term) (err error) {
	desc := strings.Join(ctx.Args[2:], " ")
	if len(desc) > 1800 {
		_, err = ctx.Sendf("Description too long (%v > 1800).", len(desc))
		return
	}

	err = c.DB.UpdateDesc(t.ID, desc)
	if err != nil {
		_, err = ctx.Sendf("Error updating description: ```%v```", err)
		return
	}

	_, err = ctx.Send("Description updated!", nil)
	if err != nil {
		c.Report(ctx, err)
	}

	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		e.Author = &discord.EmbedAuthor{
			Name: "Previous version",
		}

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title:       "Description updated",
					Description: desc,
					Color:       db.EmbedColour,
					Timestamp:   discord.NowTimestamp(),
				},
				*e,
			},
		})
	}
	return
}

func (c *Admin) editTermSource(ctx *bcr.Context, t *db.Term) (err error) {
	source := strings.Join(ctx.Args[2:], " ")
	if len(source) > 200 {
		_, err = ctx.Sendf("Source too long (%v > 200).", len(source))
		return
	}

	err = c.DB.UpdateSource(t.ID, source)
	if err != nil {
		_, err = ctx.Sendf("Error updating source: ```%v```", err)
		return
	}

	_, err = ctx.Send("Source updated!", nil)
	if err != nil {
		c.Report(ctx, err)
	}

	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		e.Author = &discord.EmbedAuthor{
			Name: "Previous version",
		}

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title:       "Source updated",
					Description: source,
					Color:       db.EmbedColour,
					Timestamp:   discord.NowTimestamp(),
				},
				*e,
			},
		})
	}
	return
}

func (c *Admin) editTermAliases(ctx *bcr.Context, t *db.Term) (err error) {
	var aliases []string
	if ctx.Args[2] != "clear" {
		aliases = strings.Split(strings.Join(ctx.Args[2:], " "), "\n")
	}

	if len(strings.Join(aliases, ", ")) > 1000 {
		_, err = ctx.Sendf("Total length of aliases too long (%v > 1000)", len(strings.Join(aliases, ", ")))
		return
	}

	err = c.DB.UpdateAliases(t.ID, aliases)
	if err != nil {
		_, err = ctx.Sendf("Error updating aliases: ```%v```", err)
		return
	}

	_, err = ctx.Send("Aliases updated!", nil)
	if err != nil {
		c.Report(ctx, err)
	}

	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		e.Author = &discord.EmbedAuthor{
			Name: "Previous version",
		}

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title:       "Aliases updated",
					Description: strings.Join(aliases, ", "),
					Color:       db.EmbedColour,
					Timestamp:   discord.NowTimestamp(),
				},
				*e,
			},
		})
	}
	return
}

func (c *Admin) editTermImage(ctx *bcr.Context, t *db.Term) (err error) {
	img := strings.Join(ctx.Args[2:], " ")
	if img == "clear" {
		img = ""
	}

	err = c.DB.UpdateImage(t.ID, img)
	if err != nil {
		_, err = ctx.Sendf("Error updating image: ```%v```", err)
		return
	}

	_, err = ctx.Send("Image updated!", nil)
	if err != nil {
		c.Report(ctx, err)
	}

	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		e.Author = &discord.EmbedAuthor{
			Name: "Previous version",
		}

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title: "Image updated",
					Image: &discord.EmbedImage{
						URL: img,
					},
					Color:     db.EmbedColour,
					Timestamp: discord.NowTimestamp(),
				},
				*e,
			},
		})
	}
	return
}

func (c *Admin) editTermTags(ctx *bcr.Context, t *db.Term) (err error) {
	var tags []string
	if ctx.Args[2] != "clear" {
		tags = ctx.Args[2:]
	}
	fmt.Println(tags)

	err = c.DB.UpdateTags(t.ID, tags)
	if err != nil {
		_, err = ctx.Sendf("Error updating tags: ```%v```", err)
		return
	}

	_, err = ctx.Send("Tags updated!", nil)
	if err != nil {
		c.Report(ctx, err)
	}

	if c.WebhookClient != nil {
		e := t.TermEmbed(c.Config.TermBaseURL())

		e.Author = &discord.EmbedAuthor{
			Name: "Previous version",
		}

		c.WebhookClient.Execute(webhook.ExecuteData{
			Username:  ctx.Bot.Username,
			AvatarURL: ctx.Bot.AvatarURL(),

			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Icon: ctx.Author.AvatarURL(),
						Name: fmt.Sprintf("%v#%v\n(%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
					},
					Title:       "Tags updated",
					Description: strings.Join(tags, ", "),
					Color:       db.EmbedColour,
					Timestamp:   discord.NowTimestamp(),
				},
				*e,
			},
		})
	}
	return
}
