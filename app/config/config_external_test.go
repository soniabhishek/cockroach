package config_test

import (
	"fmt"

	"gitlab.com/playment-main/angel/app/config"
)

func Example() {

	//set the environment
	//this has to be done only once at the app start
	config.SetEnvironment(config.Development)

	//Call GetVal to get the config value from the key
	//Keys are also exported by the package
	baseApiUrl := config.GetVal(config.BASE_API_URL)

	fmt.Println(baseApiUrl)
	// Output: localhost:8999/
}
