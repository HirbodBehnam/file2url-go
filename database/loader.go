package database

import "os"

// LoadDatabaseFromEnv will load a database from environment variables
func LoadDatabaseFromEnv() (Interface, error) {
	if path := os.Getenv("BOLT_DB_PATH"); path != "" {
		return NewBoltDatabase(path)
	} else {
		return NewMemoryCache(), nil
	}
}
