package flu_output

import "gitlab.com/playment-main/angel/app/models"

type FluOutput interface {
	AddToOutputQueue(models.FeedLineUnit) error
	AddManyToOutputQueue([]models.FeedLineUnit) error
}

func New() FluOutput {
	return &FluMonitor{}
}
