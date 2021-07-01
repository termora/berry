package pronouns

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/termora/berry/db"
)

func (bot *commands) reactionAdd(m *gateway.MessageReactionAddEvent) {
	// if this isn't the pronoun channel, return
	if m.ChannelID != bot.Config.Bot.Support.PronounChannel || !m.ChannelID.IsValid() {
		return
	}

	// this probably shouldn't happen, but check anyway
	if m.Member == nil {
		return
	}

	// if the user is a bot return
	if m.Member.User.Bot {
		return
	}

	var exists bool

	con, cancel := bot.DB.Context()
	defer cancel()

	err := bot.DB.Pool.QueryRow(con, "select exists (select * from pronoun_msgs where message_id = $1)", m.MessageID).Scan(&exists)
	if err != nil {
		bot.Sugar.Errorf("Error getting pronoun message: %v", err)
		return
	}
	if !exists {
		return
	}

	// if it's not the approve emoji, return
	if m.Emoji.Name != "✅" {
		return
	}

	var isStaff bool
	for _, r := range m.Member.RoleIDs {
		for _, s := range bot.Config.Bot.Permissions.Directors {
			if r == s {
				isStaff = true
				break
			}
		}
	}

	// if the member isn't staff, return
	if !isStaff {
		// also remove their reaction if possible
		if p, _ := bot.Router.State.Permissions(m.ChannelID, bot.Router.Bot.ID); p.Has(discord.PermissionManageMessages) && m.Emoji.Name == "✅" {
			bot.Router.State.DeleteUserReaction(m.ChannelID, m.MessageID, m.UserID, "✅")
		}

		return
	}

	var p db.PronounSet

	con, cancel = bot.DB.Context()
	defer cancel()

	err = bot.DB.Pool.QueryRow(con, "select subjective, objective, poss_det, poss_pro, reflexive from pronoun_msgs where message_id = $1", m.MessageID).Scan(&p.Subjective, &p.Objective, &p.PossDet, &p.PossPro, &p.Reflexive)
	if err != nil {
		bot.Sugar.Errorf("Error getting pronoun set: %v", err)
		return
	}

	// add the pronouns!
	_, err = bot.DB.AddPronoun(p)
	if err != nil {
		bot.Sugar.Errorf("Error adding pronoun set: %v", err)
		// this is the only one we DM the person who approved it for
		ch, chErr := bot.Router.State.CreatePrivateChannel(m.Member.User.ID)
		if chErr != nil {
			return
		}
		bot.Router.State.SendMessage(ch.ID, fmt.Sprintf("There was an error adding the pronoun set:\n```%v```", err), nil)
		return
	}

	// remove the message
	con, cancel = bot.DB.Context()
	defer cancel()

	_, err = bot.DB.Pool.Exec(con, "delete from pronoun_msgs where message_id = $1", m.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error deleting message from database: %v", err)
		return
	}

	// get the message
	msg, err := bot.Router.State.Message(m.ChannelID, m.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error getting message: %v", err)
		return
	}

	if len(msg.Embeds) < 1 {
		return
	}

	e := msg.Embeds[0]
	e.Fields = append(e.Fields, discord.EmbedField{
		Name:  "Submitted at",
		Value: fmt.Sprintf("> %v", msg.Timestamp.Time().Format("Jan 02 2006, 15:04:05 UTC")),
	})
	e.Footer = &discord.EmbedFooter{
		Icon: m.Member.User.AvatarURL(),
		Text: fmt.Sprintf("Approved by %v#%v\n(%v)", m.Member.User.Username, m.Member.User.Discriminator, m.Member.User.ID),
	}
	e.Timestamp = discord.NowTimestamp()

	_, err = bot.Router.State.EditEmbed(msg.ChannelID, msg.ID, e)
	if err != nil {
		bot.Sugar.Errorf("Error editing message: %v", err)
		return
	}
}
