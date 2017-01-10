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

	plog.Error("testing", s1{124}, plog.MessageWithParam(log_tags.FLU_ID, s2{1212, "asdokads"}), plog.Message("Asd"))
	time.Sleep(time.Duration(10) * time.Second)
}
