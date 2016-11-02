package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Mode       string
	Port       string
	SigningKey string
}

func DecodeConfiguration() Configuration {
	file, _ := os.Open("config/settings.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Fatal("Error:", err)
	}
	return conf
}
