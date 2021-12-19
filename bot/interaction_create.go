package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/getsentry/sentry-go"
	"github.com/starshine-sys/bcr"
)

// InteractionCreate is called when an interaction create event is received.
func (bot *Bot) InteractionCreate(ic *gateway.InteractionCreateEvent) {
	if ic.Data.InteractionType() != discord.CommandInteractionType {
		return
	}

	defer func() {
		r := recover()
		if r != nil {
			bot.Log.Errorf("Caught panic in channel ID %v (guild %v): %v", ic.ChannelID, ic.GuildID, r)
			bot.Log.Infof("Panicking command: %v", ic.Data.(*discord.CommandInteraction).Name)

			// if something causes a panic, it's our problem, because *it shouldn't panic*
			// so skip checking the error and just immediately report it
			var eventID *sentry.EventID
			if bot.UseSentry {
				eventID = bot.Sentry.Recover(r)
			}

			if eventID == nil {
				return
			}

			s := "An internal error has occurred. If this issue persists, please contact the bot developer with the error code above."
			if bot.Config != nil {
				if bot.Config.Bot.Support.Invite != "" {
					s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v) with the error code above.", bot.Config.Bot.Support.Invite)
				}
			}

			st, _ := bot.Router.StateFromGuildID(0)

			st.RespondInteraction(ic.ID, ic.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(fmt.Sprintf("Error code: `%v`", string(*eventID))),
					Embeds: &[]discord.Embed{{
						Title:       "Internal error occurred",
						Description: s,
						Color:       bcr.ColourRed,

						Footer: &discord.EmbedFooter{
							Text: string(*eventID),
						},
						Timestamp: discord.NowTimestamp(),
					}},
					Flags: api.EphemeralResponse,
				},
			})
		}
	}()

	ctx, err := bot.Router.NewSlashContext(ic)
	if err != nil {
		bot.Log.Errorf("Couldn't create slash context: %v", err)
		return
	}

	err = bot.Router.ExecuteSlash(ctx)
	if err != nil {
		bot.Log.Errorf("Couldn't execute slash command: %v", err)
	}

	bot.Stats.IncCommand()
}
