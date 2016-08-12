package step

import (
	"errors"

	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
)

type Step struct {
	InQ  feed_line.Fl
	OutQ feed_line.Fl

	buffer feed_line.Bf
}

func New() Step {
	return Step{
		InQ:    feed_line.New(),
		OutQ:   feed_line.New(),
		buffer: feed_line.NewBuffer(),
	}
}

func (s *Step) AddToBuffer(flu feed_line.FLU) {
	s.buffer.Save(flu)
}

func (s *Step) GetBuffered() map[uuid.UUID]feed_line.FLU {

	return s.buffer.GetAll()
}

func (s *Step) RemoveFromBuffer(flu feed_line.FLU) error {

	flu, ok := s.buffer.Get(flu.ID)
	if !ok {
		return errors.New("not present")
	}
	return nil

}
