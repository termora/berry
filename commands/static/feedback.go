package static

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (bot *Bot) feedback(ctx *bcr.Context) (err error) {
	if bot.Config.Bot.FeedbackChannel == 0 {
		_, err = ctx.Send("Sorry, but we're not currently accepting feedback through this command. Feel free to join the support server though!")
		return
	}

	for _, u := range bot.Config.Bot.FeedbackBlockedUsers {
		if u == ctx.Author.ID {
			_, err = ctx.Send("You are blocked from submitting feedback through this command. If you believe this is an error, please contact the developers.")
			return
		}
	}

	if ctx.RawArgs == "" {
		_, err = ctx.Send("You need to actually give feedback to send!")
		return
	}

	if ctx.Message.GuildID.IsValid() {
		_ = ctx.State.DeleteMessage(ctx.Message.ChannelID, ctx.Message.ID, "")
	}

	msg, err := ctx.Send("React with ✅ to send, or with ❌ to cancel.", discord.Embed{
		Description: ctx.RawArgs,
		Color:       db.EmbedColour,
	})
	if err != nil {
		return err
	}

	if yes, timeout := ctx.YesNoHandler(*msg, ctx.Author.ID); !yes || timeout {
		_, err = ctx.Send("Cancelled.")
		return
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURL(),
			Name: fmt.Sprintf("%v#%v (%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.ID),
		},
		Description: ctx.RawArgs,

		Fields: []discord.EmbedField{{
			Name:  "Source",
			Value: fmt.Sprintf("https://discord.com/channels/%v/%v/%v", ctx.Message.GuildID, ctx.Message.ChannelID, ctx.Message.ID),
		}},

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("DM from %v#%v", ctx.Author.Username, ctx.Author.Discriminator),
		},
		Timestamp: discord.NowTimestamp(),
		Color:     db.EmbedColour,
	}

	if ctx.Guild != nil {
		e.Footer.Text = fmt.Sprintf("Guild: %v (%v)", ctx.Guild.Name, ctx.Guild.ID)
	}

	_, err = ctx.State.SendEmbeds(bot.Config.Bot.FeedbackChannel, e)
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}

	if ctx.Message.GuildID.IsValid() {
		_ = ctx.State.DeleteMessage(msg.ChannelID, msg.ID, "")
		_, err = ctx.NewDM(ctx.Author.ID).Content("Thanks for submitting feedback!").Send()
	} else {
		_, err = ctx.Send("Thanks for submitting feedback!")
	}
	return
}
