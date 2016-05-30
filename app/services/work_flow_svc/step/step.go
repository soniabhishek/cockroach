package step

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

type Step struct {
	InQ  feed_line.FL
	OutQ feed_line.FL

	buffer feed_line.BF
}

func New() Step {
	return Step{
		InQ:  feed_line.New(),
		OutQ: feed_line.New(),
	}
}

func (s *Step) AddToBuffer(flu feed_line.FLU) error {

	s.buffer = append(s.buffer, flu)
	return nil
}

func (s *Step) RemoveFromBuffer(flu feed_line.FLU) error {

	return nil
}
