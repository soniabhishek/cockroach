package start_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
	//	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	//	"fmt"
	//	"reflect"
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/models"
	"github.com/pkg/errors"
)

type starterStep struct {
	step.Step
}

func (t *starterStep) processFlu(flu feed_line.FLU) {

	var c = flu.Data["image_url"]
	if c == nil || c.(string) == "" {
		t.OutQ.Push(flu)
		flu.ConfirmReceive()
		return
	}

	plog.Debug("####", flu.Data)

	t.AddToBuffer(flu)
	plog.Info("Start Step flu reached", flu.ID)

	//Image encryption
	urlSlice, err := GetEncryptedUrls(c)

	if err != nil {
		plog.Error("Image Encryption step", err, "fluId: "+flu.ID.String(), flu.FeedLineUnit)
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.StartStep, "Error in start step", flu.Redelivered())
		return
	}

	flu.Data.Merge(models.JsonF{"image_url": urlSlice})

	t.finishFlu(flu)

}

func (t *starterStep) finishFlu(flu feed_line.FLU) bool {

	err := t.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("Start step", "flu not present in buffer")
		//return false
	}
	t.OutQ.Push(flu)
	flu.ConfirmReceive()
	plog.Info("Start step out", flu.ID)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.StartStep, flu.Redelivered())
	return true
}

func GetEncryptedUrls(imageField interface{}) (urlSlice []string, err error) {

	plog.Debug("enc", imageField, "{}")
	encResult, err := clients.GetLuigiClient().GetEncryptedUrls(imageField)
	plog.Debug("ERROR", err, "{}")

	if err != nil {
		plog.Debug("slsl", encResult)
	}
	var d = 0
	for _, item := range encResult {
		if item["valid"] == false {
			err = errors.New("Image not found")
			plog.Error("Image Encryption step : Image not encryptable", err)
			return
		}
		urlSlice[d] = item["playment_url"].(string)
		d++
		plog.Debug("item", item)
	}
	return
}

/*
func test(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			fmt.Println(s.Index(i))
		}
	}
}
*/

func newStdStarter() *starterStep {
	ts := &starterStep{
		Step: step.New(step_type.StartStep),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
