package auditlog

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *AuditLog) sendPublicEmbed(e Entry, description string) (id discord.MessageID, err error) {
	if !bot.Config.Bot.AuditLogPublic.IsValid() {
		return
	}

	msg, err := bot.State.SendEmbeds(bot.Config.Bot.AuditLogPublic, bot.publicEmbed(e, description))
	if err != nil {
		return
	}

	return msg.ID, err
}

func (bot *AuditLog) publicEmbed(e Entry, description string) discord.Embed {
	embed := discord.Embed{
		Description: description,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Action ID: %v", e.ID),
		},
		Timestamp: discord.NewTimestamp(e.Timestamp),
	}

	if e.Reason.Valid {
		embed.Description += ". Reason: " + e.Reason.String
	}

	switch e.Action {
	case CreateAction:
		embed.Color = bcr.ColourGreen
	case UpdateAction:
		embed.Color = bcr.ColourBlue
	case DeleteAction:
		embed.Color = bcr.ColourRed
	}

	return embed
}

func (bot *AuditLog) sendPrivateEmbed(e Entry) (id discord.MessageID, err error) {
	if !bot.Config.Bot.AuditLogPrivate.IsValid() {
		return
	}

	msg, err := bot.State.SendEmbeds(bot.Config.Bot.AuditLogPrivate, bot.privateEmbeds(e)...)
	if err != nil {
		return
	}

	return msg.ID, err
}

func (bot *AuditLog) privateEmbeds(entry Entry) (es []discord.Embed) {
	if entry.Subject == TermEntry {
		return bot.privateTermEmbeds(entry)
	}

	name := entry.UserID.String()
	u, err := bot.State.User(entry.UserID)
	if err == nil {
		name = fmt.Sprintf("%v\n(%v)", u.Tag(), u.ID)
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{Name: name},
		Title:  fmt.Sprintf("%vd %v", strings.Title(string(entry.Action)), entry.Subject),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Action ID: %v", entry.ID),
		},
		Timestamp: discord.NewTimestamp(entry.Timestamp),
	}

	switch entry.Action {
	case CreateAction:
		e.Color = bcr.ColourGreen

		if entry.Subject == PronounsEntry {
			p, _ := entry.AfterPronouns()
			e.Description = p.String()
		} else {
			ex, _ := entry.AfterExplanation()
			e.Description = ex.Description
			e.Title += "`" + ex.Name + "`"
		}

	case UpdateAction:
		e.Color = bcr.ColourBlue

		if entry.Subject == PronounsEntry {
			before, _ := entry.BeforePronouns()
			after, _ := entry.AfterPronouns()
			e.Description = fmt.Sprintf("**%v** ➜ **%v**", before.String(), after.String())
		} else {
			ex, _ := entry.BeforeExplanation()
			e.Description = "**Before:**\n" + ex.Description
			e.Title += "`" + ex.Name + "`"
		}

	case DeleteAction:
		e.Color = bcr.ColourRed

		if entry.Subject == PronounsEntry {
			before, _ := entry.BeforePronouns()
			e.Description = before.String()
		} else {
			ex, _ := entry.BeforeExplanation()
			e.Description = "**Before:**\n" + ex.Description
			e.Title += "`" + ex.Name + "`"
		}
	}

	if entry.Reason.Valid {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Reason",
			Value: entry.Reason.String,
		})
	}

	return []discord.Embed{e}
}

func (bot *AuditLog) privateTermEmbeds(entry Entry) (es []discord.Embed) {
	before, _ := entry.BeforeTerm()
	after, _ := entry.AfterTerm()

	name := entry.UserID.String()
	u, err := bot.State.User(entry.UserID)
	if err == nil {
		name = fmt.Sprintf("%v\n(%v)", u.Tag(), u.ID)
	}

	es = []discord.Embed{{
		Author: &discord.EmbedAuthor{Name: name},
		Title:  fmt.Sprintf("%vd term \"%v\"", strings.Title(string(entry.Action)), before.Name),
		Footer: &discord.EmbedFooter{Text: fmt.Sprintf("Action ID: %v", entry.ID)},
	}}

	switch entry.Action {
	case CreateAction:
		es[0].Color = bcr.ColourGreen
	case UpdateAction:
		es[0].Color = bcr.ColourBlue
	case DeleteAction:
		es[0].Color = bcr.ColourRed
	}

	if entry.Action == CreateAction {
		if entry.Reason.Valid {
			es[0].Fields = append(es[0].Fields, discord.EmbedField{
				Name:  "Reason",
				Value: entry.Reason.String,
			})
		}

		return es
	}

	if entry.Action == DeleteAction {
		if entry.Reason.Valid {
			es[0].Fields = append(es[0].Fields, discord.EmbedField{
				Name:  "Reason",
				Value: entry.Reason.String,
			})
		}

		es = append(es, bot.DB.TermEmbed(&before))
		return es
	}

	switch {
	case before.Name != after.Name:
		es[0].Fields = append(es[0].Fields, discord.EmbedField{
			Name:  "Name",
			Value: fmt.Sprintf("**Before:** %v\n**After:** %v", before.Name, after.Name),
		})
	case strings.Join(before.Aliases, ", ") != strings.Join(after.Aliases, ", "):
		bf := "None"
		if len(before.Aliases) > 0 {
			bf = strings.Join(before.Aliases, ", ")
		}
		aft := "None"
		if len(after.Aliases) > 0 {
			strings.Join(after.Aliases, ", ")
		}

		es[0].Fields = append(es[0].Fields, []discord.EmbedField{
			{
				Name:  "Synonyms before",
				Value: bf,
			},
			{
				Name:  "Synonyms after",
				Value: aft,
			},
		}...)
	case strings.Join(before.Tags, ", ") != strings.Join(after.Tags, ", "):
		bf := "None"
		if len(before.Tags) > 0 {
			bf = strings.Join(before.Tags, ", ")
		}
		aft := "None"
		if len(after.Tags) > 0 {
			strings.Join(after.Tags, ", ")
		}

		es[0].Fields = append(es[0].Fields, []discord.EmbedField{
			{
				Name:  "Tags before",
				Value: bf,
			},
			{
				Name:  "Tags after",
				Value: aft,
			},
		}...)
	case before.ContentWarnings != after.ContentWarnings:
		bf := "None"
		if before.ContentWarnings != "" {
			bf = before.ContentWarnings
		}
		aft := "None"
		if after.ContentWarnings != "" {
			aft = after.ContentWarnings
		}

		es[0].Fields = append(es[0].Fields, []discord.EmbedField{
			{
				Name:  "CW before",
				Value: bf,
			},
			{
				Name:  "CW after",
				Value: aft,
			},
		}...)
	case before.Source != after.Source:
		es[0].Fields = append(es[0].Fields, []discord.EmbedField{
			{
				Name:  "Source before",
				Value: before.Source,
			},
			{
				Name:  "Source after",
				Value: after.Source,
			},
		}...)
	case before.Description != after.Description:
		e := discord.Embed{
			Title:       "Description updated",
			Description: before.Description,
			Color:       bcr.ColourBlue,
		}

		desc := after.Description
		if len(desc) >= 1024 {
			desc = desc[1020:] + "..."
		}
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "After",
			Value: desc,
		})

		es = append(es, e)
	}

	if entry.Reason.Valid {
		es = append(es, discord.Embed{
			Title:       "Reason",
			Description: entry.Reason.String,
		})
	}

	return es
}
