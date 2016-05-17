package util

import (
	"strings"
	"time"
)

const (
	Empty              = ""
	Star               = "*"
	Hyphen             = "-"
	UnderScore         = "_"
	WhiteSpace         = " "
	Colon              = ":"
	Dot                = "."
	Table_Referencer   = "." //Yeah, they both are same.
	Comma              = ","
	Spaced_Comma       = " , "
	Column_Quote       = "\""
	Place_Holder       = "%s"
	Place_Holder_Cover = "(%s)"

	Left_Parentheses  = "("
	Right_Parentheses = ")"
	Parentheses       = "()"
	Left_Braces       = "{"
	Right_Braces      = "{"
	Braces            = "{}"
	Left_Bracket      = "["
	Right_Bracket     = "]"
	Bracket           = "[]"
)

func GetRandomID() int64 {
	return time.Now().UnixNano()
}

func IsEmptyOrNil(str string) bool {
	return str == Empty
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
	if strings.TrimSpace(tag) == Place_Holder || strings.TrimSpace(tag) == Place_Holder_Cover {
		return true
	}
	return false
}

func areEqual_Without_Whitespaces(str1 string, str2 string) bool {
	return strings.TrimSpace(str1) == strings.TrimSpace(str2)
}
