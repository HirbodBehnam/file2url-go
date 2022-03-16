package config

import (
	"os"
	"strconv"
	"strings"
)

const Version = "1.5.0"

// Config of the application
var Config struct {
	// List of allowed users
	AllowedUsers []int64
	// Listen Address
	ListenAddress string
	// The url to show when we generate a link
	URLPrefix string
}

// LoadConfig loads the configs from environment variables
func LoadConfig() {
	// Load listen address
	if port := os.Getenv("PORT"); port != "" { // Heroku
		Config.ListenAddress = ":" + port
		Config.URLPrefix = "https://" + os.Getenv("DYNO_NAME") + ".herokuapp.com"
	} else { // Normal
		Config.ListenAddress = os.Getenv("LISTEN")
		Config.URLPrefix = os.Getenv("URL_PREFIX")
	}
	// Load allowed users
	users := strings.Split(os.Getenv("ALLOWED_USERS"), ",")
	Config.AllowedUsers = make([]int64, 0, len(users))
	for _, user := range users {
		if user, err := strconv.ParseInt(user, 10, 64); err == nil {
			Config.AllowedUsers = append(Config.AllowedUsers, user)
		}
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
