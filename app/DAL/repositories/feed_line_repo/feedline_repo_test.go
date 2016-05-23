package feed_line_repo_test

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func Example() {
	flr := feed_line_repo.New()
	flr.GetById(uuid.NewV4())
}
