package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	BotToken         string
	TestChat         string
	CobaltiumGeneral string
	CobaltiumBot     string
	RiotKey          string
}

func GetConfig() Configuration {
	f, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	config := Configuration{}

	_ = json.Unmarshal(f, &config)

	return config
}
