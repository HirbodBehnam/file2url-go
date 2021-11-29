package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config of the application
var Config struct {
	// List of allowed users
	AllowedUsers []int64 `json:"allowed_users"`
	// Listen Address
	ListenAddress string `json:"listen_address"`
	// The url to show when we generate a link
	URLPrefix string `json:"url_prefix"`
}

// LoadConfig loads the config file from a location
func LoadConfig(location string) {
	bytes, err := os.ReadFile(location)
	if err != nil {
		log.Fatalln("Cannot read config file:", err)
	}
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		log.Fatalln("Cannot parse config file:", err)
	}
}

// IsUserAllowed checks if a user is allowed to work with bot or not
func IsUserAllowed(userID int64) bool {
	// Check public bot
	if len(Config.AllowedUsers) == 0 {
		return true
	}
	// Check each allowed user
	for _, id := range Config.AllowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}
