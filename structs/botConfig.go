package structs

import "github.com/diamondburned/arikawa/v2/discord"

// FallbackGitURL if there's no git url set in the config file fall back to this
const FallbackGitURL = "https://github.com/termora/berry"

// BotConfig ...
type BotConfig struct {
	Auth struct {
		Token       string
		DatabaseURL string `json:"database_url"`
		SentryURL   string `json:"sentry_url"`
	}
	Bot struct {
		Prefixes []string

		BotOwners    []discord.UserID  `json:"bot_owners"`
		AdminServers []discord.GuildID `json:"admin_servers"`

		Support struct {
			Invite         string
			PronounChannel discord.ChannelID `json:"pronoun_channel"`
		}

		Website string
		Git     string

		// Whether to show term and server counts in the status
		ShowTermCount  bool `json:"show_term_count"`
		ShowGuildCount bool `json:"show_guild_count"`

		AllowedBots []discord.UserID `json:"allowed_bots"`

		JoinLogChannel discord.ChannelID `json:"join_log_channel"`

		TermChangelogPing string `json:"term_changelog_ping"`

		HelpFields   []discord.EmbedField `json:"help_fields"`
		CreditFields []discord.EmbedField `json:"credit_fields"`
	}

	// Fields used for sharding
	Sharded   bool `json:"-"`
	Shard     int  `json:"-"`
	NumShards int  `json:"-"`

	// UseSentry: when false, don't use Sentry for logging errors
	UseSentry bool `json:"-"`

	// Debug will print more logs
	Debug bool `json:"-"`
}

// TermBaseURL returns the base URL for terms
func (c BotConfig) TermBaseURL() string {
	if c.Bot.Website == "" {
		return ""
	}
	return c.Bot.Website + "term/"
}
