package structs

import "github.com/diamondburned/arikawa/v2/discord"

// BotConfig ...
type BotConfig struct {
	Auth struct {
		Token       string
		DatabaseURL string `json:"database_url"`
	}
	Bot struct {
		Prefixes []string

		BotOwners   []string `json:"bot_owners"`
		AdminServer string   `json:"admin_server"`

		Support struct {
			Invite         string
			SupportChannel string `json:"support_channel"`
		}

		TermBaseURL string `json:"term_base_url"`
		Website     string

		ShowGuildCount bool              `json:"show_guild_count"`
		AllowedBots    []string          `json:"allowed_bots"`
		JoinLogChannel discord.ChannelID `json:"join_log_channel"`

		TermChangelogPing string `json:"term_changelog_ping"`

		HelpFields   []EmbedField `json:"help_fields"`
		CreditFields []EmbedField `json:"credit_fields"`
	}
}

// EmbedField ...
type EmbedField struct {
	Name  string
	Value string
}
