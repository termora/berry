package bot

import (
	"context"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

// Prefixer ...
func (bot *Bot) Prefixer(m discord.Message) (n int) {
	prefixes := append(bot.Router.Prefixes, bot.PrefixesFor(m.GuildID)...)
	for _, p := range prefixes {
		if strings.HasPrefix(strings.ToLower(m.Content), p) {
			return len(p)
		}
	}
	return -1
}

// PrefixesFor returns the prefixes for the given server
func (bot *Bot) PrefixesFor(id discord.GuildID) (s []string) {
	if !id.IsValid() {
		return bot.Config.Bot.Prefixes
	}

	err := bot.DB.Pool.QueryRow(context.Background(), "select prefixes from public.servers where id = $1", id.String()).Scan(&s)
	if err != nil {
		bot.Sugar.Errorf("Error getting prefixes for %v: %v", id, err)
		// return the default prefixes
		return bot.Config.Bot.Prefixes
	}

	return s
}
