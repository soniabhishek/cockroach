package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	Development = "development_example"
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
