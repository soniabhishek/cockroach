package config

import (
	"os"

	"fmt"

	"github.com/crowdflux/angel/app/support"
	"github.com/spf13/viper"
)

const (
	Development = "development"
	Production  = "production"
	Staging     = "staging"
)

var env string
var configProvider = viper.New()

func init() {

	// Get GOENV from ~/.bash_profile or equivalent
	if env = os.Getenv("GOENV"); env == "" {
		env = Development
	}

	fmt.Println("Using " + env + " environment")

	// Derive the config directory
	configPath := support.GetAngelRootDir() + "/app/config"

	configProvider.SetConfigFile(configPath + "/" + env + ".json")
	err := configProvider.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// Returns true if current environment is Development
func IsDevelopment() bool {
	return env == Development
}

// Returns true if current environment is Production
func IsProduction() bool {
	return env == Production
}

// Returns true if current environment is Staging
func IsStaging() bool {
	return env == Staging
}

func GetEnv() string {
	return env
}
