package search

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (bot *Bot) doAutocomplete(ev *gateway.InteractionCreateEvent) {
	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	respond := func(choices []api.AutocompleteChoice) {
		_ = s.RespondInteraction(ev.ID, ev.Token, api.InteractionResponse{
			Type: api.AutocompleteResult,
			Data: &api.InteractionResponseData{
				Choices: &choices,
			},
		})
	}

	dat, ok := ev.Data.(*discord.AutocompleteInteraction)
	if !ok {
		return
	}

	if dat.Name == "explain" {
		exs, err := bot.DB.GetAllExplanations()
		if err != nil {
			bot.Log.Errorf("Error getting explanations: %v", err)
			return
		}

		choices := make([]api.AutocompleteChoice, 0, len(exs))
		for _, ex := range exs {
			choices = append(choices, api.AutocompleteChoice{Name: ex.Name, Value: ex.Name})
		}

		respond(choices)
		return
	}

	if dat.Name != "search" && dat.Name != "define" {
		return
	}

	var searchTerm string
	for _, opt := range dat.Options {
		if opt.Name == "query" {
			searchTerm = opt.Value
			break
		}
	}

	if searchTerm == "" {
		respond([]api.AutocompleteChoice{{
			Name:  "Start typing to search...",
			Value: "_this_will_not_match_anything",
		}})
	}

	terms, err := bot.DB.Autocomplete(searchTerm)
	if err != nil {
		bot.Log.Errorf("Error doing autocomplete search: %v", err)
		return
	}

	opts := make([]api.AutocompleteChoice, 0, len(terms))
	for _, t := range terms {
		opts = append(opts, api.AutocompleteChoice{Name: t, Value: t})
	}

	respond(opts)
}
