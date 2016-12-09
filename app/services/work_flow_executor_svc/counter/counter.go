package counter

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/plog"
)

func Print(flu feed_line.FLU, step string) {
	plog.Info(step, "reached: "+flu.ID.String())
}
