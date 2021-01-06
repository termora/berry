package structs

// BotConfig ...
type BotConfig struct {
	Auth struct {
		Token       string
		DatabaseURL string `json:"database_url"`
	}
	Bot struct {
		Prefixes     []string
		BotOwners    []string `json:"bot_owners"`
		AdminServer  string   `json:"admin_server"`
		ServerInvite string   `json:"server_invite"`
		Website      string
		TermBaseURL  string   `json:"term_base_url"`
		AllowedBots  []string `json:"allowed_bots"`

		HelpFields []HelpField `json:"help_fields"`
	}
}

// HelpField ...
type HelpField struct {
	Name  string
	Value string
}
