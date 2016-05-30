package feed_line

import (
	"gitlab.com/playment-main/angel/app/models"
)

// ShortHand for channel of FLUs i.e. FeedLine
type FL chan models.FeedLineUnit

// Get new FeedLine channel with unlimited size
func New() FL {

	feedLine := make(chan models.FeedLineUnit, 1000)
	return feedLine
}

func NewBig() FL {
	feedLine := make(chan models.FeedLineUnit, 10000)
	return feedLine
}

// Get new FeedLine channel with fixed size
func NewFixedSize(size int) FL {
	return make(chan models.FeedLineUnit, size)
}

//--------------------------------------------------------------------------------//

type BF []models.FeedLineUnit
