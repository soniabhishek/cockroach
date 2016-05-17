package feed_line_repo

import "gitlab.com/playment-main/support/app/services/plerrors"

var ErrDuplicateReferenceId plerrors.ServiceError = plerrors.ServiceError{"FLR_0001", "Duplicate Reference Id"}
var ErrFLUNotFoundInInputQueue plerrors.ServiceError = plerrors.ServiceError{"FLR_0002", "FeedLineUnit not present in input queue"}
