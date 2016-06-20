package manual_step

import (
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

type manualStep struct {
	step.Step
}

func (m *manualStep) processFlu(flu feed_line.FLU) {
	m.AddToBuffer(flu)
	plog.Info("Manual Step flu reached", flu)
}

//TODO call it once upload method gets called.
func (m *manualStep) finishFlu(flu feed_line.FLU) bool {

	err := m.RemoveFromBuffer(flu)
	if err != nil {
		return false
	}
	counter.Print(flu, "manual")
	m.OutQ <- flu
	return true
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
