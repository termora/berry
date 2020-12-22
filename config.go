package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/Starshine113/termbot/structs"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func getConfig(sugar *zap.SugaredLogger) (config *structs.BotConfig) {
	token := flag.String("token", "", "Override the token in config.yaml")
	databaseURL := flag.String("db", "", "Override the database URL in config.yaml")
	flag.Parse()

	config = &structs.BotConfig{}

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		sampleConf, err := ioutil.ReadFile("config.sample.yaml")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("config.yaml", sampleConf, 0644)
		if err != nil {
			panic(err)
		}
		sugar.Errorf("config.yaml was not found, created sample configuration.")
		os.Exit(1)
		return nil
	}
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &config)
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
