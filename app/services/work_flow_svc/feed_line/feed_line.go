package feed_line

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

// ShortHand for channel of FLUs i.e. FeedLine
type FL chan FLU

// Get new FeedLine channel with unlimited size
func New() FL {

	feedLine := make(chan FLU, 1000)
	return feedLine
}

func NewBig() FL {
	feedLine := make(chan FLU, 10000)
	return feedLine
}

// Get new FeedLine channel with fixed size
func NewFixedSize(size int) FL {
	return make(chan FLU, size)
}

//--------------------------------------------------------------------------------//

type BF []FLU

//--------------------------------------------------------------------------------//

type FLU struct {
	models.FeedLineUnit

	// Change the name
	// Used here for collecting & passing around information about
	// previous steps of the flu processing
	Trip []Builder
}

//--------------------------------------------------------------------------------//

type Builder interface {
	GetStep() step.StepIdentifier
	GetData() interface{}
}
