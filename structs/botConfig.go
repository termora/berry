package structs

import "github.com/diamondburned/arikawa/v2/discord"

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
			SupportChannel string `json:"support_channel"`
		}

		TermBaseURL string `json:"term_base_url"`
		Website     string

		ShowGuildCount bool              `json:"show_guild_count"`
		AllowedBots    []discord.UserID  `json:"allowed_bots"`
		JoinLogChannel discord.ChannelID `json:"join_log_channel"`

		TermChangelogPing string `json:"term_changelog_ping"`

		HelpFields   []EmbedField `json:"help_fields"`
		CreditFields []EmbedField `json:"credit_fields"`
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

// EmbedField ...
type EmbedField struct {
	Name  string
	Value string
}
