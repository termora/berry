package common

import "github.com/diamondburned/arikawa/v3/discord"

// FallbackGitURL if there's no git url set in the config file fall back to this
const FallbackGitURL = "https://github.com/termora/berry"

// Webhook ...
type Webhook struct {
	ID    discord.WebhookID `json:"id" toml:"id"`
	Token string            `json:"token" toml:"token"`
}

// BotConfig ...
type BotConfig struct {
	Auth struct {
		Token       string
		DatabaseURL string `json:"database_url" toml:"database_url"`
		SentryURL   string `json:"sentry_url" toml:"sentry_url"`

		TypesenseURL string `json:"typesense_url" toml:"typesense_url"`
		TypesenseKey string `json:"typesense_key" toml:"typesense_key"`

		InfluxDB struct {
			URL    string `toml:"url"`
			Token  string `toml:"token"`
			Org    string `toml:"org"`
			Bucket string `toml:"bucket"`
		} `toml:"influxdb"`
	}
	Bot struct {
		LicenseLink string `json:"license_link" toml:"license_link"`

		Prefixes []string

		BotOwners   []discord.UserID `json:"bot_owners" toml:"bot_owners"`
		Permissions struct {
			Admins    []discord.RoleID `json:"admins" toml:"admins"`
			Directors []discord.RoleID `json:"directors" toml:"directors"`
		} `json:"permissions" toml:"permissions"`

		SlashCommands struct {
			Enabled bool              `json:"enabled" toml:"enabled"`
			Guilds  []discord.GuildID `json:"guilds" toml:"guilds"` // empty to sync all guilds
		} `json:"slash_commands" toml:"slash_commands"`

		Support struct {
			Invite         string
			PronounChannel discord.ChannelID `json:"pronoun_channel" toml:"pronoun_channel"`

			// These should be the support server, and a token for a bot *in* said support server, with the guild members intent (and in the future, message content) enabled. Blame Discord.
			GuildID discord.GuildID `toml:"guild_id"`
			Token   string          `toml:"token"`
		}

		// mostly for debugging, send a webhook message when the bot shuts down
		StartStopLog Webhook `json:"start_log" toml:"start_log"`

		AuditLog struct {
			Public  discord.ChannelID `json:"public" toml:"public"`
			Private discord.ChannelID `json:"private" toml:"private"`
		} `json:"audit_log" toml:"audit_log"`

		Website string
		Git     string

		// Whether to show term and server counts in the status
		ShowTermCount  bool `json:"show_term_count" toml:"show_term_count"`
		ShowGuildCount bool `json:"show_guild_count" toml:"show_guild_count"`
		// Whether to show shard number in the status
		ShowShard bool `json:"show_shard" toml:"show_shard"`

		AllowedBots []discord.UserID `json:"allowed_bots" toml:"allowed_bots"`

		JoinLogChannel discord.ChannelID `json:"join_log_channel" toml:"join_log_channel"`

		TermChangelogPing string `json:"term_changelog_ping" toml:"term_changelog_ping"`

		HelpFields   []discord.EmbedField `json:"help_fields" toml:"help_fields"`
		CreditFields []discord.EmbedField `json:"credit_fields" toml:"credit_fields"`

		// this will be used by t;invite and t;about if set
		CustomInvite string `json:"custom_invite" toml:"custom_invite"`

		FeedbackChannel      discord.ChannelID `json:"feedback_channel" toml:"feedback_channel"`
		FeedbackBlockedUsers []discord.UserID  `json:"feedback_blocked_users" toml:"feedback_blocked_users"`
	}

	// BotLists is tokens for the two bot lists the bot is on
	// will POST guild count every hour
	BotLists struct {
		TopGG  string `json:"top.gg" toml:"topgg"`
		BotsGG string `json:"bots.gg" toml:"botsgg"`
	} `json:"bot_lists" toml:"bot_lists"`

	// QuickNotes is a map of notes that can quickly be set with `t;admin setnote`
	QuickNotes map[string]string `json:"quick_notes" toml:"quick_notes"`

	RPCPort          string            `json:"rpc_port" toml:"rpc_port"`
	ContributorRoles []ContributorRole `json:"contributor_roles" toml:"contributor_roles"`

	// UseSentry: when false, don't use Sentry for logging errors
	UseSentry bool `json:"-" toml:"-"`
}

// ContributorRole ...
type ContributorRole struct {
	Name string         `json:"name" toml:"name"`
	ID   discord.RoleID `json:"id" toml:"id"`
}

// TermBaseURL returns the base URL for terms
func (c BotConfig) TermBaseURL() string {
	if c.Bot.Website == "" {
		return ""
	}
	return c.Bot.Website + "term/"
}
