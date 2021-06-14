package structs

import "github.com/diamondburned/arikawa/v2/discord"

// FallbackGitURL if there's no git url set in the config file fall back to this
const FallbackGitURL = "https://github.com/termora/berry"

// BotConfig ...
type BotConfig struct {
	// Fields used for sharding
	NumShards int `json:"num_shards"`
	// These are not filled in the config file
	Sharded bool `json:"-"`
	Shard   int  `json:"-"`

	Auth struct {
		Token       string
		DatabaseURL string `json:"database_url"`
		SentryURL   string `json:"sentry_url"`
	}
	Bot struct {
		Prefixes []string

		BotOwners   []discord.UserID `json:"bot_owners"`
		Permissions struct {
			Admins    []discord.RoleID `json:"admins"`
			Directors []discord.RoleID `json:"directors"`
		} `json:"permissions"`

		Support struct {
			Invite         string
			PronounChannel discord.ChannelID `json:"pronoun_channel"`
		}

		TermLog struct {
			ID    discord.WebhookID
			Token string
		} `json:"term_log"`

		Website string
		Git     string

		// Whether to show term and server counts in the status
		ShowTermCount  bool `json:"show_term_count"`
		ShowGuildCount bool `json:"show_guild_count"`
		// Whether to show shard number in the status
		ShowShard bool `json:"show_shard"`

		AllowedBots []discord.UserID `json:"allowed_bots"`

		JoinLogChannel discord.ChannelID `json:"join_log_channel"`

		TermChangelogPing string `json:"term_changelog_ping"`

		HelpFields   []discord.EmbedField `json:"help_fields"`
		CreditFields []discord.EmbedField `json:"credit_fields"`

		// this will be used by t;invite and t;about if set
		CustomInvite string `json:"custom_invite"`

		FeedbackChannel      discord.ChannelID `json:"feedback_channel"`
		FeedbackBlockedUsers []discord.UserID  `json:"feedback_blocked_users"`
	}

	// BotLists is tokens for the two bot lists the bot is on
	// will POST guild count every hour
	BotLists struct {
		TopGG  string `json:"top.gg"`
		BotsGG string `json:"bots.gg"`
	} `json:"bot_lists"`

	// QuickNotes is a map of notes that can quickly be set with `t;admin setnote`
	QuickNotes map[string]string `json:"quick_notes"`

	// UseSentry: when false, don't use Sentry for logging errors
	UseSentry bool `json:"-"`
}

// TermBaseURL returns the base URL for terms
func (c BotConfig) TermBaseURL() string {
	if c.Bot.Website == "" {
		return ""
	}
	return c.Bot.Website + "term/"
}
