package structs

// BotConfig ...
type BotConfig struct {
	Auth struct {
		Token       string
		DatabaseURL string `json:"database_url"`
	}
	Bot struct {
		Prefixes       []string
		BotOwners      []string `json:"bot_owners"`
		AdminServer    string   `json:"admin_server"`
		ServerInvite   string   `json:"server_invite"`
		Website        string
		ShowGuildCount bool     `json:"show_guild_count"`
		TermBaseURL    string   `json:"term_base_url"`
		AllowedBots    []string `json:"allowed_bots"`

		TermChangelogPing string `json:"term_changelog_ping"`

		HelpFields []EmbedField `json:"help_fields"`

		CreditFields []EmbedField `json:"credit_fields"`
	}
}

// EmbedField ...
type EmbedField struct {
	Name  string
	Value string
}
