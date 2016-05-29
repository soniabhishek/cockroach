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
}

func (c *crowdSourcingStep) processFlu(flu models.FeedLineUnit) {

	qBody := someIn.DeriveQuestionBody(flu)

	question := models.Question{
		Body:     qBody,
		Label:    flu.ReferenceId + "sdf",
		IsTest:   sql.NullBool{false, true},
		IsActive: sql.NullBool{false, true},
	}

	someIn.SaveQuestion(question)

	c.AddToBuffer(flu)

	c.finishFlu(flu)
}

func (c *crowdSourcingStep) finishFlu(flu models.FeedLineUnit) {

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
