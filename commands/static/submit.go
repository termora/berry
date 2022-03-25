package static

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/common/log"
	"github.com/termora/berry/db"
)

func (bot *Bot) submitFeedback(v bcr.Contexter) (err error) {
	ctx, ok := v.(*bcr.SlashContext)
	if !ok {
		return nil
	}

	if bot.Config.Bot.FeedbackChannel == 0 {
		return ctx.SendEphemeral("Sorry, but we're not currently accepting feedback through this command. Feel free to join the support server, though!")
	}

	for _, u := range bot.Config.Bot.FeedbackBlockedUsers {
		if u == ctx.Author.ID {
			return ctx.SendEphemeral("You are blocked from submitting feedback through this command. If you believe this is an error, please contact the developers.")
		}
	}

	return ctx.State.RespondInteraction(ctx.InteractionID, ctx.InteractionToken, api.InteractionResponse{
		Type: api.ModalResponse,
		Data: &api.InteractionResponseData{
			Title:    option.NewNullableString("Submit feedback"),
			CustomID: option.NewNullableString("submit-feedback-modal"),
			Components: &discord.ContainerComponents{
				&discord.ActionRowComponent{
					&discord.TextInputComponent{
						CustomID:    "feedback",
						Style:       discord.TextInputParagraphStyle,
						Label:       "Your feedback",
						ValueLimits: [2]int{1, 4000},
						Required:    true,
					},
				},
			},
		},
	})
}

func (bot *Bot) interactionCreate(ic *gateway.InteractionCreateEvent) {
	data, ok := ic.Data.(*discord.ModalInteraction)
	if !ok {
		return
	}

	var err error
	switch data.CustomID {
	case "submit-feedback-modal":
		err = bot.handleFeedback(ic, data)
	case "submit-term-modal":
	case "submit-pronouns-modal":
		err = bot.handlePronouns(ic, data)
	}
	if err != nil {
		log.Errorf("handling modal interaction: %v", err)
	}
}

func (bot *Bot) handleFeedback(ic *gateway.InteractionCreateEvent, data *discord.ModalInteraction) (err error) {
	var feedback string
	for _, cc := range data.Components {
		v, ok := cc.(*discord.ActionRowComponent)
		if ok {
			for _, c := range *v {
				v, ok := c.(*discord.TextInputComponent)
				if ok && v.CustomID == "feedback" {
					feedback = v.Value.Val
				}
			}
		}
	}

	if feedback == "" {
		return bot.respondEphemeral(ic, "You didn't give any feedback! This is a bug.")
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ic.Sender().AvatarURL(),
			Name: fmt.Sprintf("%v (%v)", ic.Sender().Tag(), ic.Sender().ID),
		},
		Description: feedback,

		Footer: &discord.EmbedFooter{
			Text: "From /submit feedback command",
		},
		Timestamp: discord.NowTimestamp(),
		Color:     db.EmbedColour,
	}

	s, _ := bot.Router.StateFromGuildID(ic.GuildID)
	_, err = s.SendEmbeds(bot.Config.Bot.FeedbackChannel, e)
	if err != nil {
		log.Errorf("sending feedback message: %v", err)
		return bot.respondEphemeral(ic, "There was an unknown error while sending your feedback. Try again?")
	}
	return bot.respondEphemeral(ic, "Thanks for submitting feedback!")
}

func (bot *Bot) respondEphemeral(ic *gateway.InteractionCreateEvent, tmpl string, v ...interface{}) error {
	s, _ := bot.Router.StateFromGuildID(ic.GuildID)
	return s.RespondInteraction(ic.ID, ic.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf(tmpl, v...)),
			Flags:   api.EphemeralResponse,
		},
	})
}
