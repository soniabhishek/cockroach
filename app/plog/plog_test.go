package plog

import (
	"testing"
	"time"
)

type s1 struct {
	Id int
}

func (s1) Error() string {
	return "fuck up"
}

func TestErrorMail(t *testing.T) {

	t.SkipNow()

	Error("testing", s1{124})
	Error("testing", s1{125}, NewMessage("with message"))
	Error("testing", s1{124}, NewMessageWithParam("with message & args", 123), NewMessage("Asd"))
	time.Sleep(time.Duration(10) * time.Second)
}
