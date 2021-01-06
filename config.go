package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/Starshine113/berry/structs"
	"go.uber.org/zap"
)

func getConfig(sugar *zap.SugaredLogger) (config *structs.BotConfig) {
	token := flag.String("token", "", "Override the token in config.json")
	databaseURL := flag.String("db", "", "Override the database URL in config.json")
	flag.Parse()

	config = &structs.BotConfig{}

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
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, &config)
	sugar.Infof("Loaded configuration file.")

	if *token != "" {
		config.Auth.Token = *token
	}
	if *databaseURL != "" {
		config.Auth.DatabaseURL = *databaseURL
	}
	if os.Getenv("TERMBOT_DB_URL") != "" {
		config.Auth.DatabaseURL = os.Getenv("TERMBOT_DB_URL")
	}

	return config
}
