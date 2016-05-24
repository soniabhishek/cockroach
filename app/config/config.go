package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
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
	env = os.Getenv("GOENV")

	// Get goPath
	goPath := strings.Split(os.Getenv("GOPATH"), ":")[0]

	// Derive the config directory
	configPath := goPath + "/src/gitlab.com/playment-main/angel/app/config"

	configProvider.SetConfigFile(configPath + "/" + env + ".json")
	configProvider.SetConfigName(env)

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
