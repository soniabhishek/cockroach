package bulk_processor

import (
	"errors"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
)

type Job struct {
	do func()
}

func NewJob(do func()) Job {
	return Job{
		do: do,
	}
}

func (j *Job) Do() {
	defer func() {
		if r := recover(); r != nil {
			plog.Error("Job", errors.New("Panic occured in a Job"), plog.MP(log_tags.RECOVER, r))
		}
	}()

	j.do()
}

type jobChannel chan Job
