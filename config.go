package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/starshine-sys/berry/structs"
	"go.uber.org/zap"
)

func getConfig(sugar *zap.SugaredLogger) (config *structs.BotConfig) {
	config = &structs.BotConfig{}
	fn := "config.json"
	if os.Getenv("TERMBOT_CONFIG") != "" {
		fn = os.Getenv("TERMBOT_CONFIG")
	}

	if fn == "config.json" {
		if _, err := os.Stat("config.json"); os.IsNotExist(err) {
			sampleConf, err := ioutil.ReadFile("config.sample.json")
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile("config.json", sampleConf, 0644)
			if err != nil {
				panic(err)
			}
			sugar.Errorf("config.json was not found, created sample configuration.")
			os.Exit(1)
			return nil
		}
	}
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, &config)
	sugar.Infof("Loaded configuration file.")

	if os.Getenv("TERMBOT_DB_URL") != "" {
		config.Auth.DatabaseURL = os.Getenv("TERMBOT_DB_URL")
	}
	config.UseSentry = config.Auth.SentryURL != ""

	return config
}
