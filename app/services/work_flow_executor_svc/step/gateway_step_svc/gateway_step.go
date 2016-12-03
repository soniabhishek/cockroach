package gateway_step_svc

// Futuristic stuff

//
//import (
//	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
//	"github.com/crowdflux/angel/app/models/step_type"
//	"github.com/crowdflux/angel/app/plog"
//	"github.com/crowdflux/angel/app/DAL/feed_line"
//	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
//)
//
//type gateWayStep struct {
//	step.Step
//	stepRepo step_repo.IStepRepo
//}
//
//func (m *gateWayStep) processFlu(flu feed_line.FLU) {
//
//	plog.Info("Gateway Step flu reached", flu)
//
//	step, err := m.stepRepo.GetStartStep(flu.ProjectId)
//	if err != nil {
//		plog.Error("Gateway", err, flu)
//	}
//	if step.Type == step_type.Gateway {
//		flu.StepId = step.ID
//	} else {
//		//handle error
//	}
//	m.AddToBuffer(flu)
//	m.finishFlu(flu)
//}
//
//func (m *gateWayStep) finishFlu(flu feed_line.FLU) bool {
//
//	err := m.RemoveFromBuffer(flu)
//	if err != nil {
//		return false
//	}
//	m.OutQ <- flu
//	return true
//}
//
//func (m *gateWayStep) start() {
//	go func() {
//		for {
//			select {
//			case flu := <-m.InQ:
//				m.processFlu(flu)
//			}
//		}
//	}()
//}
//
//func (m *gateWayStep) Connect(routerIn *feed_line.Fl) (routerOut *feed_line.Fl) {
//
//	// Send output of this step to the router's input
//	// for next rerouting
//	m.OutQ = *routerIn
//
//	m.start()
//
//	// Return the input channel of this step
//	// so that router can push flu to it
//	return &m.InQ
//}
