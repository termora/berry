package pronouns

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/termora/berry/db"
)

func (c *commands) use(ctx bcr.Contexter) (err error) {
	pronouns := ctx.GetStringFlag("pronouns")
	name := ctx.GetStringFlag("name")
	if v, ok := ctx.(*bcr.Context); ok {
		if len(v.Args) == 0 {
			return ctx.SendEphemeral(
				fmt.Sprintf("You didn't give any pronouns to show! Try ``%vlist-pronouns`` for a list of all pronouns.", c.Config.Bot.Prefixes[0]))
		}

		pronouns = v.Args[0]
		if len(v.Args) > 1 {
			name = v.Args[1]
		}
	}

	if pronouns == "" {
		return ctx.SendEphemeral(
			fmt.Sprintf("You didn't give any pronouns to show! Try ``%vlist-pronouns`` for a list of all pronouns.", c.Config.Bot.Prefixes[0]))
	}

	sets, err := c.DB.GetPronoun(strings.Split(pronouns, "/")...)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return ctx.SendEphemeral(
				fmt.Sprintf("Couldn't find any pronoun sets from your input. Try `%vlist-pronouns` for a list of all pronouns; if it's not on there, feel free to submit it with `%vsubmit-pronouns`!", c.Config.Bot.Prefixes[0], c.Config.Bot.Prefixes[0]))
		}
		if err == db.ErrTooManyForms {
			return ctx.SendEphemeral("You gave too many forms! Input up to five forms, separated with a slash (`/`).")
		}
		return c.DB.InternalError(ctx, err)
	}

	if len(sets) > 1 {
		if len(sets) > 25 {
			return ctx.SendEphemeral("Found more than 25 sets matching your input! Please try again.")
		}
		return c.pronounList(ctx, sets, name)
	}
	// use the first set
	set := sets[0]

	if tmplCount == 0 {
		return ctx.SendEphemeral("There are no examples available for pronouns! If you think this is in error, please join the bot support server and ask there.")
	}

	useSet := &db.PronounSet{
		Subjective: set.Subjective,
		Objective:  set.Objective,
		PossDet:    set.PossDet,
		PossPro:    set.PossPro,
		Reflexive:  set.Reflexive,
	}
	if name != "" {
		useSet.Subjective = name
	}

	e, err := c.pronounEmbeds(set, useSet)

	if v, ok := ctx.(*bcr.Context); ok {
		_, err = v.PagedEmbed(e, false)
	} else {
		_, _, err = ctx.ButtonPages(e, 15*time.Minute)
	}
	return
}

func (c *commands) pronounEmbeds(set, useSet *db.PronounSet) (e []discord.Embed, err error) {
	var b strings.Builder

	e = append(e, discord.Embed{
		Title:       fmt.Sprintf("%v/%v pronouns", set.Subjective, set.Objective),
		Description: fmt.Sprintf("**%s**\n\nTo see these pronouns in action, use the arrow reactions on this message!", set),
		Color:       db.EmbedColour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Page 1/%v", set.ID, tmplCount+1),
		},
	})

	for i := 0; i < tmplCount; i++ {
		err = templates.ExecuteTemplate(&b, strconv.Itoa(i), useSet)
		if err != nil {
			return
		}
		e = append(e, discord.Embed{
			Title:       fmt.Sprintf("%v/%v pronouns", set.Subjective, set.Objective),
			Description: b.String(),
			Color:       db.EmbedColour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("ID: %v | Page %v/%v", set.ID, i+2, tmplCount+1),
			},
		})
		b.Reset()
	}

	return e, err
}

func (c *commands) pronounList(ctx bcr.Contexter, sets []*db.PronounSet, name string) (err error) {
	s := fmt.Sprintf("Found more than one set matching your input! Please select the set you want to use:")

	options := []discord.SelectComponentOption{}

	for i, set := range sets {
		options = append(options, discord.SelectComponentOption{
			Label: set.String(),
			Value: fmt.Sprint(i),
		})
	}

	comp := []discord.Component{&discord.ActionRowComponent{Components: []discord.Component{&discord.SelectComponent{
		CustomID:    "pronouns",
		Options:     options,
		Placeholder: "Select a pronoun set...",
	}}}}

	msg, err := ctx.SendComponents(comp, s)
	if err != nil {
		return
	}

	con, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ignoreFn := func(ev *gateway.InteractionCreateEvent) bool {
		components := discord.UnwrapComponents(ev.Message.Components)

		err := ctx.Session().RespondInteraction(ev.ID, ev.Token, api.InteractionResponse{
			Type: api.UpdateMessage,
			Data: &api.InteractionResponseData{
				Components: &components,
			},
		})
		if err != nil {
			c.Sugar.Errorf("Error responding to interaction: %v", err)
		}

		return false
	}

	var ind int
	v := ctx.Session().WaitFor(con, func(v interface{}) bool {
		ev, ok := v.(*gateway.InteractionCreateEvent)
		if !ok {
			return false
		}

		if ev.Data == nil || ev.Message == nil {
			return false
		}

		data, ok := ev.Data.(*discord.ComponentInteractionData)
		if !ok {
			return false
		}

		if data.ComponentType != discord.SelectComponentType || ev.Message.ID != msg.ID {
			return false
		}

		u := ev.User
		if ev.User == nil {
			u = &ev.Member.User
		}

		if u.ID != ctx.User().ID {
			return ignoreFn(ev)
		}

		ind, err = strconv.Atoi(data.Values[0])
		if err != nil {
			return ignoreFn(ev)
		}

		return true
	})

	comp = []discord.Component{&discord.ActionRowComponent{Components: []discord.Component{&discord.SelectComponent{
		CustomID:    "pronouns",
		Options:     options,
		Placeholder: "Select a pronoun set...",
		Disabled:    true,
	}}}}

	ctx.EditOriginal(api.EditInteractionResponseData{
		Components: &comp,
	})

	if v == nil {
		return
	}

	set := sets[ind]
	useSet := &db.PronounSet{
		Subjective: set.Subjective,
		Objective:  set.Objective,
		PossDet:    set.PossDet,
		PossPro:    set.PossPro,
		Reflexive:  set.Reflexive,
	}
	if name != "" {
		useSet.Subjective = name
	}

	e, err := c.pronounEmbeds(set, useSet)
	if err != nil {
		return c.DB.InternalError(ctx, err)
	}

	ev := v.(*gateway.InteractionCreateEvent)

	// replace interaction ID/token with new one
	if v, ok := ctx.(*bcr.SlashContext); ok {
		v.InteractionID = ev.ID
		v.InteractionToken = ev.Token

		_, _, err = v.ButtonPages(e, 15*time.Minute)
		if err != nil {
			return err
		}
	} else {
		err := ctx.Session().RespondInteraction(ev.ID, ev.Token, api.InteractionResponse{
			Type: api.UpdateMessage,
			Data: &api.InteractionResponseData{
				Components: &comp,
			},
		})
		if err != nil {
			c.Sugar.Errorf("Error responding to interaction: %v", err)
		}

		_, _, err = ctx.ButtonPages(e, 15*time.Minute)
		if err != nil {
			return err
		}
	}

	err = ctx.Session().DeleteMessage(msg.ChannelID, msg.ID, "")
	if err != nil {
		c.Sugar.Errorf("Error deleting message: %v", err)
	}
	return nil
}
