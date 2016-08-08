package main

import (
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
)

func main() {

	feed_line_repo.SyncAll()

}
