package feed_line

import (
	"gitlab.com/playment-main/angel/app/models"
)

// ShortHand for channel of FLUs i.e. FeedLine
type Fl chan FLU

// Get new FeedLine channel with unlimited size
func New() Fl {

	feedLine := make(chan FLU, 1000)
	return feedLine
}

func NewBig() Fl {
	feedLine := make(chan FLU, 10000)
	return feedLine
}

// Get new FeedLine channel with fixed size
func NewFixedSize(size int) Fl {
	return make(chan FLU, size)
}

//--------------------------------------------------------------------------------//

type Bf []FLU

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
	GetStep() uint
	GetData() interface{}
}
