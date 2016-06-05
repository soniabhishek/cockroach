package step

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

type Step struct {
	InQ  feed_line.Fl
	OutQ feed_line.Fl

	buffer feed_line.Bf
}

func New() Step {
	return Step{
		InQ:  feed_line.New(),
		OutQ: feed_line.New(),
	}
}

func (s *Step) AddToBuffer(flu feed_line.FLU) {
	s.buffer.Add(flu)
}

func (s *Step) GetBuffered() feed_line.Bf {

	return s.buffer
}

func (s *Step) RemoveFromBuffer(flu feed_line.FLU) error {

	return nil
}

func (s *Step) Detain(flu feed_line.FLU, why error, saver iFluSave) {
	err := saver.Save(flu.FeedLineUnit)
	if err != nil {
		plog.Error("Step error", err)
	}

}

type iFluSave interface {
	Save(models.FeedLineUnit)
}
