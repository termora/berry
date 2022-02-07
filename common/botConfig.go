package common

import "github.com/diamondburned/arikawa/v3/discord"

// FallbackGitURL if there's no git url set in the config file fall back to this
const FallbackGitURL = "https://github.com/termora/berry"

// Webhook ...
type Webhook struct {
	ID    discord.WebhookID `toml:"id"`
	Token string            `toml:"token"`
}

// BotConfig ...
type BotConfig struct {
	Token    string `toml:"token"`
	InfluxDB struct {
		URL    string `toml:"url"`
		Token  string `toml:"token"`
		Org    string `toml:"org"`
		Bucket string `toml:"bucket"`
	} `toml:"influxdb"`

	LicenseLink string `toml:"license_link"`

	Prefixes []string

	BotOwners []discord.UserID `toml:"bot_owners"`
	Admins    []discord.RoleID `toml:"admins"`
	Directors []discord.RoleID `toml:"directors"`

	SlashEnabled bool              `toml:"slash_commands_enabled"`
	SlashGuilds  []discord.GuildID `toml:"slash_commands_guilds"` // empty to sync all guilds

	SupportInvite  string            `toml:"support_invite"`
	PronounChannel discord.ChannelID `toml:"pronoun_channel"`

	// These should be the support server, and a token for a bot *in* said support server, with the guild members intent (and in the future, message content) enabled. Blame Discord.
	SupportGuildID discord.GuildID `toml:"support_guild_id"`
	SupportToken   string          `toml:"support_token"`

	// mostly for debugging, send a webhook message when the bot shuts down
	StartStopLog Webhook `toml:"start_log"`

	AuditLogPublic  discord.ChannelID `toml:"audit_log_public"`
	AuditLogPrivate discord.ChannelID `toml:"audit_log_private"`

	Website string

	// Whether to show term and server counts in the status
	ShowTermCount  bool `toml:"show_term_count"`
	ShowGuildCount bool `toml:"show_guild_count"`
	// Whether to show shard number in the status
	ShowShard bool `toml:"show_shard"`

	AllowedBots []discord.UserID `toml:"allowed_bots"`

	JoinLogChannel discord.ChannelID `toml:"join_log_channel"`

	TermChangelogPing string `toml:"term_changelog_ping"`

	// this will be used by t;invite and t;about if set
	CustomInvite string `toml:"custom_invite"`

	FeedbackChannel      discord.ChannelID `toml:"feedback_channel"`
	FeedbackBlockedUsers []discord.UserID  `toml:"feedback_blocked_users"`

	// BotLists is tokens for the two bot lists the bot is on
	// will POST guild count every hour
	TopGG  string `toml:"topgg"`
	BotsGG string `toml:"botsgg"`

	HelpFields   []discord.EmbedField `toml:"help_fields"`
	CreditFields []discord.EmbedField `toml:"credit_fields"`

	// QuickNotes is a map of notes that can quickly be set with `t;admin setnote`
	QuickNotes map[string]string `toml:"quick_notes"`

	ContributorRoles []ContributorRole `toml:"contributor_roles"`
}

// ContributorRole ...
type ContributorRole struct {
	Name string         `toml:"name"`
	ID   discord.RoleID `toml:"id"`
}

// TermBaseURL returns the base URL for terms
func (c BotConfig) TermBaseURL() string {
	if c.Website == "" {
		return ""
	}
	return c.Website + "term/"
}
