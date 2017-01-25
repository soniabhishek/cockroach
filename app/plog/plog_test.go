package plog_test

import (
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
	"testing"
	"time"
)

type s1 struct {
	Id int
}

type s2 struct {
	Id int
	pc string
}

func (s1) Error() string {
	return "fuck up"
}

func TestErrorMail(t *testing.T) {

	t.SkipNow()
	plog.Error("testing", s1{124}, plog.MessageWithParam(log_tags.FLU_ID, s2{1212, "asdokads"}), plog.Message("Asd"))
	plog.Error("testing", s1{124})
	plog.Error("testing", s1{125}, plog.M("with message"))
	plog.Error("testing", s1{124}, plog.MessageWithParam(log_tags.ROW_NUM, 123))
	time.Sleep(time.Duration(10) * time.Second)
}

func TestCustomLogger(t *testing.T) {
	c := plog.NewLogger("somelog", "INFO", "FILE")

	c.Info("testingjhoh", s1{124}, "")
	c.Info("testihihjing", s1{125}, "with message")
	c.Info("teshhiting", s1{124}, "with message & args", 123, "Asd")
	time.Sleep(time.Duration(10) * time.Second)
}
