package counter

import (
	"fmt"

	"github.com/crowdflux/angel/app/DAL/feed_line"
)

func Print(flu feed_line.FLU, step string) {
	fmt.Println(flu.ID.String() + " reached in " + step)
}
