package common

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/termora/berry/common/log"
)

func ReadConfig() Config {
	var config Config

	fn := "config.toml"
	if os.Getenv("TERMORA_CONFIG") != "" {
		fn = os.Getenv("TERMORA_CONFIG")
	}

	configFile, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("Couldn't find or open file: %v", err)
	}

	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Couldn't unmarshal config file: %v", err)
	}

	log.Infof("Loaded configuration file %q.", fn)

	// get some config keys from environment, for Docker
	if s := os.Getenv("DATABASE_URL"); s != "" {
		config.Core.DatabaseURL = s
	}
	if s := os.Getenv("TYPESENSE_URL"); s != "" {
		config.Core.TypesenseURL = s
	}
	if s := os.Getenv("TYPESENSE_APIKEY"); s != "" {
		config.Core.TypesenseKey = s
	}
	if s := os.Getenv("REDIS"); s != "" {
		config.Core.Redis = s
	}
	if s := os.Getenv("PORT"); s != "" {
		// these are always run as separate processes, so it's easier to just have both use the same env var
		config.Site.Port = s
		config.API.Port = s
	}

	config.Core.UseSentry = config.Core.SentryURL != ""

	if config.Core.Git == "" {
		config.Core.Git = FallbackGitURL
	}
	config.Site.Git = config.Core.Git

	return config
}

type Config struct {
	Core CoreConfig `toml:"core"`
	Bot  BotConfig  `toml:"bot"`
	Site SiteConfig `toml:"site"`
	API  APIConfig  `toml:"api"`
}

type CoreConfig struct {
	DatabaseURL string `toml:"database_url"`
	SentryURL   string `toml:"sentry_url"`

	TypesenseURL string `toml:"typesense_url"`
	TypesenseKey string `toml:"typesense_key"`

	Git string `toml:"git"`

	Redis string `toml:"redis"` // optional

	// UseSentry: when false, don't use Sentry for logging errors
	UseSentry bool `toml:"-"`
}

type SiteConfig struct {
	Port string `toml:"port"`

	SiteName string `toml:"site_name"`
	BaseURL  string `toml:"base_url"`
	Invite   string `toml:"invite"`
	Contact  bool   `toml:"contact"`
	// Optional description shown in embeds, when not linking to a term page
	Description string `toml:"description"`

	Plausible struct {
		Domain string `toml:"domain"`
		URL    string `toml:"url"`
	} `toml:"plausible"`

	Git string `toml:"-"`
}

type APIConfig struct {
	Port string `toml:"port"`
}
