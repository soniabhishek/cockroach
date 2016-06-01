package crowdsourcing_step

import (
	"database/sql"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

type crowdSourcingStep struct {
	step.Step
	fluManager fluManager
}

func (c *crowdSourcingStep) processFlu(flu feed_line.FLU) {

	qBody := c.fluManager.DeriveQuestionBodyFromFlu(flu)

	question := models.Question{
		Body:     qBody,
		Label:    flu.ReferenceId + "sdf",
		IsTest:   sql.NullBool{false, true},
		IsActive: sql.NullBool{false, true},
	}

	err := c.fluManager.QuestionRepo.Add(question)
	if err != nil {
		//c.Detain(flu, err, c.fluManager.QuestionRepo)
		return
	}

	c.AddToBuffer(flu)

}

func (c *crowdSourcingStep) finishFlu(flu feed_line.FLU) {

	c.RemoveFromBuffer(flu)
	flu.Step = "crowdsourcing"
	counter.Print(flu)
	c.OutQ <- flu
}

func (c *crowdSourcingStep) start() {
	go func() {
		for {
			select {
			case flu := <-c.InQ:
				c.processFlu(flu)
			}
		}
	}()
}

func (c *crowdSourcingStep) Connect(routerIn *feed_line.Fl) (routerOut *feed_line.Fl) {

	// Send output of this step to the router's input
	// for next rerouting
	c.OutQ = *routerIn

	c.start()

	// Return the input channel of this step
	// so that router can push flu to it
	return &c.InQ
}

func (c *crowdSourcingStep) HandleQuestionComplete(q models.Question, ans models.QuestionAnswer) error {
	panic("Not implemented")

	//flu := c.fluManager.FeedlineRepo.GetByQuestionId(q.ID)

	//c.finishFlu(flu)
	return nil
}
