package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/termora/berry/structs"
	"go.uber.org/zap"
)

func getConfig(sugar *zap.SugaredLogger) structs.BotConfig {
	var config structs.BotConfig

	fn := "config"
	if os.Getenv("TERMBOT_CONFIG") != "" {
		fn = os.Getenv("TERMBOT_CONFIG")
	}

	fullName := fn + ".toml"
	format := "toml"
	if _, err := os.Stat(fn + ".toml"); os.IsNotExist(err) {
		fullName = fn + ".json"
		format = "json"
	}
	configFile, err := ioutil.ReadFile(fullName)
	if err != nil {
		sugar.Fatalf("Couldn't find or open file: %v", err)
	}
	switch format {
	case "toml":
		err = toml.Unmarshal(configFile, &config)
	case "json":
		err = json.Unmarshal(configFile, &config)
	}
	if err != nil {
		sugar.Fatalf("Couldn't unmarshal config file: %v", err)
	}

	sugar.Infof("Loaded configuration file \"%v\".", fullName)

	if os.Getenv("TERMBOT_DB_URL") != "" {
		config.Auth.DatabaseURL = os.Getenv("TERMBOT_DB_URL")
	}
	config.UseSentry = config.Auth.SentryURL != ""

	if config.Bot.Git == "" {
		config.Bot.Git = structs.FallbackGitURL
	}

	return config
}
