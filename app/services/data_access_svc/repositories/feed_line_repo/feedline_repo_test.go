package feed_line_repo_test

import (
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/data_access_svc/repositories/feed_line_repo"
)

func Example() {
	flr := feed_line_repo.New()
	flr.GetById(uuid.NewV4())
}
