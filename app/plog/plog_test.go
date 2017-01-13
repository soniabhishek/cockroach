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
	Error("testing", s1{124}, "")
	Error("testing", s1{125}, "with message")
	Error("testing", s1{124}, "with message & args", 123, "Asd")
	time.Sleep(time.Duration(10) * time.Second)
}

func TestCustomLogger(t *testing.T) {
	c := NewLogger("somelog", "INFO", "FILE")

	c.Info("testingjhoh", s1{124}, "")
	c.Info("testihihjing", s1{125}, "with message")
	c.Info("teshhiting", s1{124}, "with message & args", 123, "Asd")
	time.Sleep(time.Duration(10) * time.Second)
}
