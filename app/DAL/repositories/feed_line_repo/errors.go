package feed_line_repo

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/plerrors"
)

var ErrDuplicateReferenceId plerrors.ServiceError = plerrors.ServiceError{"FLR_0001", "Duplicate Reference Id"}
var ErrFLUNotFoundInInputQueue plerrors.ServiceError = plerrors.ServiceError{"FLR_0002", "FeedLineUnit not present in input queue"}
var ErrFLUNotFound plerrors.ServiceError = plerrors.ServiceError{"FLR_0003", "FeedLineUnit not found"}
var ErrNoUpdatableFlus plerrors.ServiceError = plerrors.ServiceError{"FLR_0004", "No Flus to update"}
var ErrPartiallyUpdatedFlus plerrors.ServiceError = plerrors.ServiceError{"FLR_0005", "Partially Update the flus"}

type BulkInsertError struct {
	Error     error
	BulkError []BulkError
}

type BulkError struct {
	Message string
	Flu     models.FeedLineUnit
}
