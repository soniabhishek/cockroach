package crowdsourcing_step

import (
	"database/sql"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/question_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

type crowdSourcingStep struct {
	step.Step
	questionRepo question_repo.IQuestionRepo
	fluRepo      feed_line_repo.IFluRepo
}

func (c *crowdSourcingStep) processFlu(flu feed_line.FLU) {

	qBody := someIn.DeriveQuestionBody(flu.FeedLineUnit)

	question := models.Question{
		Body:     qBody,
		Label:    flu.ReferenceId + "sdf",
		IsTest:   sql.NullBool{false, true},
		IsActive: sql.NullBool{false, true},
	}

	err := c.questionRepo.Add(question)
	if err != nil {

	}

	c.AddToBuffer(flu)

	c.finishFlu(flu)
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

func (c *crowdSourcingStep) Connect(routerIn *feed_line.FL) (routerOut *feed_line.FL) {

	// Send output of this step to the router's input
	// for next rerouting
	c.OutQ = *routerIn

	c.start()

	// Return the input channel of this step
	// so that router can push flu to it
	return &c.InQ
}

//--------------------------------------------------------------------------------//

type testInterface interface {
	DeriveQuestionBody(models.FeedLineUnit) models.JsonFake
	SaveQuestion(models.Question)
}

var someIn testInterface = testStruct{}

type testStruct struct {
}

func (testStruct) DeriveQuestionBody(flu models.FeedLineUnit) models.JsonFake {
	jf := models.JsonFake{
		"body": "1234",
		"flu":  flu,
	}
	//fmt.Println("questionderived", flu.ID)
	return jf
}

func (testStruct) SaveQuestion(q models.Question) {
	//fmt.Println("question saved")
}
