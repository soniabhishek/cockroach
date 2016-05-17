/*
This package is for managing the configuration keys & environments.
	config.GetVal(config.BASE_API_URL)
*/
package config

import (
	"sync"
)

type configKeyValues map[configKey]string

var singleTonConfig configKeyValues

var once sync.Once

/*
Returns the string value of the given input key.

Key is of type config.configKey.

Type config.configKey isn't exported (for obvious reasons).

All the exported keys can be found at environment.go
*/
func GetVal(c configKey) string {

	once.Do(initConfiguration)

	if val, ok := singleTonConfig[c]; ok {
		return val
	}

	panic(newConfigurationError(c))

}

/*
This function will get executed only once
since it uses sync.Once internally
*/
func initConfiguration() {

	switch singleTonEnv {
	case Production:
		singleTonConfig = getProductionConfiguration()
	case Development:
		singleTonConfig = getDevelopmentConfiguration()
	case Staging:
		singleTonConfig = getStagingConfiguration()
	default:
		singleTonConfig = getDevelopmentConfiguration()
	}
}
