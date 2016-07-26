package transformation_step

import (
	"strings"

	"gitlab.com/playment-main/angel/app/DAL/repositories/projects_repo"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

type transformationStep struct {
	step.Step
	projectsRepo projects_repo.IProjectsRepo
}

func (m *transformationStep) processFlu(flu feed_line.FLU) {
	m.AddToBuffer(flu)
	plog.Info("transformation Step flu reached", flu)

	proj, err := m.projectsRepo.GetById(flu.ProjectId)
	if err != nil {
		plog.Error("Transformation", err, "error loading project")
		return
	}

	if strings.Contains(strings.ToLower(proj.Name), "flip") {
		flipkartHack(flu)
	}

}

func (m *transformationStep) finishFlu(flu feed_line.FLU) bool {

	err := m.RemoveFromBuffer(flu)
	if err != nil {
		return false
	}
	m.OutQ <- flu
	plog.Info("transformation out", flu.ID)
	return true
}

func (m *transformationStep) start() {
	go func() {
		for {
			select {
			case flu := <-m.InQ:
				m.processFlu(flu)
			}
		}
	}()
}

func (m *transformationStep) Connect(routerIn *feed_line.Fl) (routerOut *feed_line.Fl) {

	// Send output of this step to the router's input
	// for next rerouting
	m.OutQ = *routerIn

	m.start()

	// Return the input channel of this step
	// so that router can push flu to it
	return &m.InQ
}
