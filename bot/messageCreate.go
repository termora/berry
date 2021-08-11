package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/getsentry/sentry-go"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

// MessageCreate is run when a message is created and handles commands
func (bot *Bot) MessageCreate(m *gateway.MessageCreateEvent) {
	var err error
	var ctx *bcr.Context

	// defer panic handling
	defer func() {
		r := recover()
		if r != nil {
			bot.Sugar.Errorf("Caught panic in channel ID %v (user %v, guild %v): %v", m.ChannelID, m.Author.ID, m.GuildID, r)
			bot.Sugar.Infof("Panic message content:\n```\n%v\n```", m.Content)

			// if something causes a panic, it's our problem, because *it shouldn't panic*
			// so skip checking the error and just immediately report it
			var eventID *sentry.EventID
			if bot.UseSentry {
				eventID = bot.Sentry.Recover(r)
			}

			if ctx == nil || eventID == nil {
				return
			}

			s := "An internal error has occurred. If this issue persists, please contact the bot developer with the error code above."
			if bot.Config != nil {
				if bot.Config.Bot.Support.Invite != "" {
					s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v) with the error code above.", bot.Config.Bot.Support.Invite)
				}
			}

			ctx.Send(
				fmt.Sprintf("Error code: `%v`", string(*eventID)),
				discord.Embed{
					Title:       "Internal error occurred",
					Description: s,
					Color:       bcr.ColourRed,

					Footer: &discord.EmbedFooter{
						Text: string(*eventID),
					},
					Timestamp: discord.NowTimestamp(),
				},
			)
		}
	}()

	// if the bot user isn't set yet, do that here
	// we can't do it when initialising the router because the connection to Discord will error
	if bot.Router.Bot == nil {
		err = bot.Router.SetBotUser(m.GuildID)
		if err != nil {
			bot.Sugar.Error("Error setting bot user:", err)
			return
		}
		bot.Router.Prefixes = append(bot.Router.Prefixes, fmt.Sprintf("<@%v>", bot.Router.Bot.ID), fmt.Sprintf("<@!%v>", bot.Router.Bot.ID))
	}

	// if message was sent by a bot return, unless it's in the list of allowed bots
	if m.Author.Bot && !inSlice(bot.Config.Bot.AllowedBots, m.Author.ID) {
		return
	}
	// if the message content is empty (indicating an embed-only bot message), return
	if m.Content == "" {
		return
	}

	// get context
	ctx, err = bot.Router.NewContext(m)
	if err != nil && err != bcr.ErrEmptyMessage {
		bot.Sugar.Error("Error creating context:", err)
		return
	}

	// check if the message might be a command
	if bot.Router.MatchPrefix(m.Message) {
		bot.Sugar.Debugf("Maybe executing command `%v`", ctx.Command)

		err = bot.Router.Execute(ctx)
		if err != nil {
			if db.IsOurProblem(err) && bot.UseSentry {
				bot.DB.CaptureError(ctx, err)
			}
			bot.Sugar.Error(err)
		}
	}
}

func inSlice(slice []discord.UserID, s discord.UserID) bool {
	for _, i := range slice {
		if i == s {
			return true
		}
	}
	return false
}
