package config

import "sync"

type environment int32

const (
	Development environment = iota
	Staging
	Production
)

var singleTonEnv environment
var once2 sync.Once

//This function will get executed only once
//as it uses sync.Once internally
func SetEnvironment(env environment) {
	once2.Do(func() {
		singleTonEnv = env
	})
}

func GetEnvironment() environment {
	return singleTonEnv
}
