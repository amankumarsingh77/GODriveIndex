package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	Auth   AuthConfig
	UI     UIConfig
	Player PlayerConfig
)

func LoadConfig() {
	loadJSONConfig("auth.json", &Auth)
	// loadJSONConfig("config/ui.json", &UI)
	// loadJSONConfig("config/player.json", &Player)
}

func loadJSONConfig(filename string, config interface{}) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading config file %s: %v", filename, err)
	}
	if err := json.Unmarshal(data, config); err != nil {
		log.Fatalf("Error parsing config file %s: %v", filename, err)
	}
}
