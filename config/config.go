package config

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration matches the fields in settings.json
type Configuration struct {
	Mode       string
	Port       string
	SigningKey string
}

// DecodeConfiguration gets the settings data from settings.json
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
