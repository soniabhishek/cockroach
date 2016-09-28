package config_temp

import (
	"github.com/spf13/viper"
)

var Vipe = viper.New()

func init() {
	Vipe.SetConfigFile("/Users/himanshu144141/code/gocode/src/github.com/crowdflux/angel/config_temp/test.json")
	err := Vipe.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
