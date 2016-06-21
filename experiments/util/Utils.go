package util

import (
	"gitlab.com/playment-main/angel/utilities/constants"
	"strings"
	"time"
)

func GetRandomID() int64 {
	return time.Now().UnixNano()
}

func IsEmptyOrNil(str string) bool {
	return str == constants.Empty
}
func IsValidError(err error) bool {
	return err != nil
}

func IsValidTag(tag string) bool {
	if IsEmptyOrNil(strings.TrimSpace(tag)) || isKeyword(tag) || IsPlaceHolder(tag) {
		return false
	}
	return true
}
func IsPlaceHolder(tag string) bool {
	if strings.TrimSpace(tag) == constants.Place_Holder || strings.TrimSpace(tag) == constants.Place_Holder_Cover {
		return true
	}
	return false
}

func areEqual_Without_Whitespaces(str1 string, str2 string) bool {
	return strings.TrimSpace(str1) == strings.TrimSpace(str2)
}
