package structs

// BotConfig ...
type BotConfig struct {
	Auth struct {
		Token       string
		DatabaseURL string `yaml:"database_url"`
	}
	Bot struct {
		Prefixes     []string
		BotOwners    []string `yaml:"bot_owners"`
		AdminServer  string   `yaml:"admin_server"`
		ServerInvite string   `yaml:"support_server"`
		Website      string
	}
}
