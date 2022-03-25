package static

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
)

func (bot *Bot) submitPronouns(v bcr.Contexter) (err error) {
	ctx, ok := v.(*bcr.SlashContext)
	if !ok {
		return nil
	}

	if bot.Config.Bot.PronounChannel == 0 {
		return ctx.SendEphemeral("We aren't accepting new pronoun submissions through the bot. You might be able to ask in the support server.")
	}

	return ctx.State.RespondInteraction(ctx.InteractionID, ctx.InteractionToken, api.InteractionResponse{
		Type: api.ModalResponse,
		Data: &api.InteractionResponseData{
			Title:    option.NewNullableString("Submit pronouns"),
			CustomID: option.NewNullableString("submit-pronouns-modal"),
			Components: &discord.ContainerComponents{
				&discord.ActionRowComponent{
					&discord.TextInputComponent{
						CustomID:    "subjective",
						Style:       discord.TextInputShortStyle,
						Label:       "Subjective form (example: she)",
						Placeholder: option.NewNullableString("Example: she"),
						ValueLimits: [2]int{1, 100},
						Required:    true,
					},
				},
				&discord.ActionRowComponent{
					&discord.TextInputComponent{
						CustomID:    "objective",
						Style:       discord.TextInputShortStyle,
						Label:       "Objective form (example: her)",
						Placeholder: option.NewNullableString("Example: her"),
						ValueLimits: [2]int{1, 100},
						Required:    true,
					},
				},
				&discord.ActionRowComponent{
					&discord.TextInputComponent{
						CustomID:    "poss_det",
						Style:       discord.TextInputShortStyle,
						Label:       "Possessive determiner (example: *her* pen)",
						Placeholder: option.NewNullableString("Example: *her* pen"),
						ValueLimits: [2]int{1, 100},
						Required:    true,
					},
				},
				&discord.ActionRowComponent{
					&discord.TextInputComponent{
						CustomID:    "poss_pro",
						Style:       discord.TextInputShortStyle,
						Label:       "Possessive pronoun (example: hers)",
						Placeholder: option.NewNullableString("Example: that pen is *hers*"),
						ValueLimits: [2]int{1, 100},
						Required:    true,
					},
				},
				&discord.ActionRowComponent{
					&discord.TextInputComponent{
						CustomID:    "reflexive",
						Style:       discord.TextInputShortStyle,
						Label:       "Reflexive form (example: herself)",
						Placeholder: option.NewNullableString("Example: herself"),
						ValueLimits: [2]int{1, 100},
						Required:    true,
					},
				},
			},
		},
	})
}

func (bot *Bot) handlePronouns(ic *gateway.InteractionCreateEvent, data *discord.ModalInteraction) (err error) {
	var p db.PronounSet
	for _, cc := range data.Components {
		v, ok := cc.(*discord.ActionRowComponent)
		if ok {
			for _, c := range *v {
				v, ok := c.(*discord.TextInputComponent)
				if !ok {
					continue
				}

				switch v.CustomID {
				case "subjective":
					p.Subjective = strings.ToLower(strings.TrimSpace(v.Value.Val))
				case "objective":
					p.Objective = strings.ToLower(strings.TrimSpace(v.Value.Val))
				case "poss_det":
					p.PossDet = strings.ToLower(strings.TrimSpace(v.Value.Val))
				case "poss_pro":
					p.PossPro = strings.ToLower(strings.TrimSpace(v.Value.Val))
				case "reflexive":
					p.Reflexive = strings.ToLower(strings.TrimSpace(v.Value.Val))
				}
			}
		}
	}

	if p.Subjective == "" || p.Objective == "" || p.PossDet == "" || p.PossPro == "" || p.Reflexive == "" {
		return bot.respondEphemeral(ic, "One or more required forms was empty! This is a bug.")
	}

	_, err = bot.DB.GetPronoun(p.Subjective, p.Objective, p.PossDet, p.PossPro, p.Reflexive)
	if err == nil {
		return bot.respondEphemeral(ic, "That pronoun set already exists!")
	}

	found := false
	err = bot.DB.QueryRow(context.Background(), `select exists(
		select * from pronoun_msgs where
		subjective = $1 and objective = $2
		and poss_det = $3 and poss_pro = $4
		and reflexive = $5)`, p.Subjective, p.Objective, p.PossDet, p.PossPro, p.Reflexive).Scan(&found)
	if err != nil {
		log.Errorf("error checking if pronoun set exists: %v", err)
	}

	if found {
		return bot.respondEphemeral(ic, "That pronoun set has already been submitted!")
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: fmt.Sprintf("%v (%v)", ic.Sender().Tag(), ic.SenderID()),
			Icon: ic.Sender().AvatarURL(),
		},
		Color:       db.EmbedColour,
		Title:       "Pronoun submission",
		Description: p.String(),
		Fields: []discord.EmbedField{{
			Name:  "Submitted by",
			Value: ic.Sender().Mention(),
		}},
		Timestamp: discord.NowTimestamp(),
	}

	s, _ := bot.Router.StateFromGuildID(ic.GuildID)
	msg, err := s.SendEmbeds(bot.Config.Bot.PronounChannel, e)
	if err != nil {
		log.Errorf("sending pronouns message: %v", err)
		return bot.respondEphemeral(ic, "There was an unknown error while submitting these pronouns. Try again?")
	}

	_, err = bot.DB.Exec(context.Background(), `insert into pronoun_msgs
	(message_id, subjective, objective, poss_det, poss_pro, reflexive)
	values ($1, $2, $3, $4, $5, $6)`, msg.ID, p.Subjective, p.Objective, p.PossDet, p.PossPro, p.Reflexive)
	if err == nil {
		// if the error's non-nil, the message was still sent
		// so don't just return immediately
		s.React(msg.ChannelID, msg.ID, "✅")
	} else {
		log.Errorf("Error adding submission message %v to database: %v", msg.ID, err)
	}

	return bot.respondEphemeral(ic, "Successfully submitted the pronoun set **%v**!", p.String())
}
