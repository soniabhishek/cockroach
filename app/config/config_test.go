package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Check if both development & production config list
//return the same set of keys
func TestConfigKeys(t *testing.T) {

	t.SkipNow()

	developmentConfigValues := getDevelopmentConfiguration()
	productionConfigValues := getProductionConfiguration()

	developmentKeys := getKeys(developmentConfigValues)
	productionKeys := getKeys(productionConfigValues)

	//union of development & production keys
	allKeys := append(developmentKeys, productionKeys...)

	//get uniques from allKeys
	uniqueKeys := getUniqueKeys(allKeys)

	for _, key := range uniqueKeys {

		//check if key is present in development config file
		_, ok := developmentConfigValues[key]

		assert.True(t, ok, newConfigurationError(key).Error()+" Development Config")

		//check if key is present in production config file
		_, ok = productionConfigValues[key]
		assert.True(t, ok, newConfigurationError(key).Error()+" Production Config")

	}

	//assert.Fail(t, "asd")
}

//Get array of keys from configuration key value list
func getKeys(c configKeyValues) []configKey {

	ch := make([]configKey, 0, len(c))

	for configKey := range c {
		ch = append(ch, configKey)
	}

	return ch
}

//remove duplicate config keys in a list of configKeys
func getUniqueKeys(input []configKey) []configKey {

	encountered := map[configKey]bool{}
	result := []configKey{}

	for index, key := range input {
		if encountered[key] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[key] = true
			// Append to result slice.
			result = append(result, input[index])
		}
	}
	return result
}
