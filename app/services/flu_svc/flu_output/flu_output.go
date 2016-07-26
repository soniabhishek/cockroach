package flu_output

import "github.com/crowdflux/angel/app/models"

type FluOutput interface {
	AddToOutputQueue(models.FeedLineUnit) error
	AddManyToOutputQueue([]models.FeedLineUnit) error
}

func New() FluOutput {
	StartFluOutputTimer()
	return &FluMonitor{}
}
