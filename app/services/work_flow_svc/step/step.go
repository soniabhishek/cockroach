package step

import (
	"gitlab.com/playment-main/angel/app/models"
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

func (s *Step) AddToBuffer(flu feed_line.FLU) error {

	s.buffer = append(s.buffer, flu)
	return nil
}

func (s *Step) RemoveFromBuffer(flu feed_line.FLU) error {

	return nil
}

func (s *Step) Detain(flu feed_line.FLU, why error, saver iFluSave) {
	_ = saver.Save(flu.FeedLineUnit)

}

type iFluSave interface {
	Save(models.FeedLineUnit) error
}
