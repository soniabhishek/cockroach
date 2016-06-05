package manual_step

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

type manualStep struct {
	step.Step
}

func (m *manualStep) processFlu(flu feed_line.FLU) {
	m.AddToBuffer(flu)
}

func (m *manualStep) start() {
	go func() {
		for {
			select {
			case flu := <-m.InQ:
				m.processFlu(flu)
			}
		}
	}()
}

func (m *manualStep) Connect(routerIn *feed_line.Fl) (routerOut *feed_line.Fl) {

	// Send output of this step to the router's input
	// for next rerouting
	m.OutQ = *routerIn

	m.start()

	// Return the input channel of this step
	// so that router can push flu to it
	return &m.InQ
}
