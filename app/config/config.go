package config

import (
	"os"
	"runtime"
	"strings"

	"fmt"
)

const (
	Development = "development"
	Production  = "production"
	Staging     = "staging"
)

var env string

//var configProvider = viper.New()

func init() {

	// Get GOENV from ~/.bash_profile or equivalent
	if env = os.Getenv("GOENV"); env == "" {
		env = Development
	}

	fmt.Println("Using " + env + " environment")

	var goPath string

	// windows has ; separator vs linux has :
	if runtime.GOOS == "windows" {
		goPath = strings.Split(os.Getenv("GOPATH"), ";")[0]
	} else {
		goPath = strings.Split(os.Getenv("GOPATH"), ":")[0]
	}

	_ = goPath

	//configProvider.SetConfigFile(configPath + "/" + env + ".json")
	//configProvider.SetConfigName(env)
	//
	//err := configProvider.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}
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
