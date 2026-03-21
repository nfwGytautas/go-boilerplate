// Package env provides a simple wrapper for loading and checking the
// environment variables
package env

import (
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

var (
	// EnvironmentKey defines the key used to check for the environment i.e. functions like IsProduction
	EnvironmentKey = "ENV"
	// VersionKey defines the key used to check for the version of the software
	VersionKey = "VERSION"
	// ProductionEnvValue defines the value that constitutes a production build
	ProductionEnvValue = "production"
)

// Load the environment variables (it will load a .env file if it exists)
func Load() {
	var err error

	_, err = os.Stat(".env")
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic("Load: " + err.Error())
	}

	err = godotenv.Load(".env")
	if err != nil {
		panic("Load: " + err.Error())
	}
}

// Parse calls Load and then parses the environment variables into a struct
func Parse[T any]() T {
	var result T
	err := env.Parse(&result)
	if err != nil {
		panic("Parse: " + err.Error())
	}
	return result
}

// IsProduction checks if the current EnvironmentKey variable is equal to ProductionEnvValue
func IsProduction() bool {
	val := os.Getenv(EnvironmentKey)
	if val == "" {
		panic("IsProduction used without specifying a value (check 'EnvironmentKey' and make sure you defined the value)")
	}
	return val == ProductionEnvValue
}

// Version returns the version of the software
func Version() string {
	val := os.Getenv(VersionKey)
	if val == "" {
		panic("Version used without specifying a value (check 'VersionKey' and make sure you defined the value)")
	}
	return val
}
