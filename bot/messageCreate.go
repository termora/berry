package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/getsentry/sentry-go"
	"github.com/mediocregopher/radix/v4"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/common/log"
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
			log.Errorf("Caught panic in channel ID %v (user %v, guild %v): %v", m.ChannelID, m.Author.ID, m.GuildID, r)
			log.Infof("Panic message content:\n```\n%v\n```", m.Content)

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
			if bot.Config.Bot.SupportInvite != "" {
				s = fmt.Sprintf("An internal error has occurred. If this issue persists, please contact the bot developer in the [support server](%v) with the error code above.", bot.Config.Bot.SupportInvite)
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
	if err != nil {
		if err != bcr.ErrEmptyMessage {
			log.Error("Error creating context:", err)
		}
		return
	}

	// apparently we sometimes panic on line 83, not sure what happens there--just gonna return here
	if ctx == nil {
		log.Errorf("Error was %v, but Context is nil.", err)
		return
	}

	// check if the message might be a command
	if bot.Router.MatchPrefix(m.Message) {
		log.Debugf("Maybe executing command `%v`", ctx.Command)

		err = bot.Router.Execute(ctx)
		if err != nil {
			if db.IsOurProblem(err) && bot.UseSentry {
				bot.DB.CaptureError(ctx, err)
			}
			log.Error(err)
		}

		bot.slashReminderMessage(ctx)
		bot.Stats.IncCommand()
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

func (bot *Bot) slashReminderMessage(ctx *bcr.Context) {
	if bot.redis == nil || strings.HasPrefix(ctx.Message.Content, "<@") {
		return
	}

	var i int
	err := bot.redis.Do(context.Background(), radix.Cmd(&i, "SISMEMBER", "termora:slash-reminders", ctx.Author.ID.String()))
	if err != nil {
		log.Error("redis error:", err)
		return
	}
	if i == 1 {
		return
	}

	e := discord.Embed{
		Description: fmt.Sprintf(`**Note:** text prefixes (such as `+"`%v`"+`) will no longer be supported as of April 31st.
Please use mentions (%v) or slash commands instead.
This message will only show up once.`, ctx.Prefix, ctx.Bot.Mention()),
		Color: db.EmbedColour,
	}

	_, err = ctx.SendComponents(discord.Components(&discord.ButtonComponent{
		Style:    discord.SecondaryButtonStyle(),
		CustomID: "delete-reminder",
		Label:    "Got it, delete this message",
	}), "", e)
	if err != nil {
		log.Error("error sending prefix reminder message:", err)
		return
	}

	err = bot.redis.Do(context.Background(), radix.Cmd(&i, "SADD", "termora:slash-reminders", ctx.Author.ID.String()))
	if err != nil {
		log.Error("redis error:", err)
		return
	}
}

func (bot *Bot) reminderInteraction(ic *gateway.InteractionCreateEvent) {
	data, ok := ic.Data.(*discord.ButtonInteraction)
	if !ok {
		return
	}

	if data.CustomID == "delete-reminder" && ic.Message != nil {
		s, _ := bot.Router.StateFromGuildID(ic.GuildID)
		_ = s.DeleteMessage(ic.Message.ChannelID, ic.Message.ID, "")
	}
}
