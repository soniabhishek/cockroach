package step

import (
	"errors"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"sync"
)

type Step struct {
	InQ  feed_line.Fl
	OutQ feed_line.Fl

	buffer feed_line.Bf
	once   sync.Once

	processFlu processFlu
	stepType   step_type.StepType
}

type processFlu func(feed_line.FLU)

func New(st step_type.StepType) Step {
	return Step{
		InQ:      feed_line.New(st.String() + "-in"),
		OutQ:     feed_line.New(st.String() + "-out"),
		buffer:   feed_line.NewBuffer(),
		stepType: st,
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

func (s *Step) Connect(routerIn *feed_line.Fl) *feed_line.Fl {

	// Send output of this step to the router's input
	// for next rerouting
	s.OutQ = *routerIn

	s.start()

	// Return the input channel of this step
	// so that router can push flu to it
	return &s.InQ
}

// TODO shit code
func (s *Step) SetFluProcessor(p processFlu) {
	s.processFlu = p
}

func (s *Step) start() {

	if s.processFlu == nil {
		panic(errors.New("processFlu nil for the step"))
	}

	s.once.Do(func() {

		go func() {

			for flu := range s.InQ.Receiver() {

				flu_logger_svc.LogStepEntry(flu.FeedLineUnit, s.stepType, flu.Redelivered())
				s.processFlu(flu)

			}
		}()
	})
}
