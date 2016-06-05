package counter

import (
	"fmt"

	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

func Print(flu feed_line.FLU, step string) {
	fmt.Println(flu.ID.String() + " reached in " + step)
}
