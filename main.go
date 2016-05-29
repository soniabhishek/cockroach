package main

import (
	"gitlab.com/playment-main/angel/app"
)

func main() {

	//if false {
	//
	//	defer func() {
	//		if err := recover(); err != nil {
	//			if configError, ok := err.(*config.ConfigNotFoundError); ok {
	//				fmt.Println(configError.Error())
	//			} else {
	//				panic(err)
	//			}
	//		}
	//	}()
	//}

	app.Start()
}
