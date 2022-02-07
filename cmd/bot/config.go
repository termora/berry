package bot

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/termora/berry/common"
	"github.com/termora/berry/common/log"
)

func getConfig() common.BotConfig {
	var config common.BotConfig

	fn := "config.bot"
	if os.Getenv("TERMBOT_CONFIG") != "" {
		fn = os.Getenv("TERMBOT_CONFIG")
	}

	fullName := fn + ".toml"
	configFile, err := ioutil.ReadFile(fullName)
	if err != nil {
		log.Fatalf("Couldn't find or open file: %v", err)
	}

	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Couldn't unmarshal config file: %v", err)
	}

	log.Infof("Loaded configuration file \"%v\".", fullName)

	if os.Getenv("TERMBOT_DB_URL") != "" {
		config.Auth.DatabaseURL = os.Getenv("TERMBOT_DB_URL")
	}
	config.UseSentry = config.Auth.SentryURL != ""

	if config.Bot.Git == "" {
		config.Bot.Git = common.FallbackGitURL
	}

	return config
}
