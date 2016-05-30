package utilities

import (
	"fmt"
	"time"

	"github.com/beefsack/go-rate"
)

var rl = rate.New(1000, time.Hour) //TODO configurable

func IsItUnderLimit(message string) bool {
	if ok, remaining := rl.Try(); ok {
		fmt.Println("Remaining time: ", remaining)
		return true
	} else {
		return false
	}
}
