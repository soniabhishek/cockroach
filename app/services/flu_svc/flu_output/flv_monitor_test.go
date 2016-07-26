package flu_output

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	//"github.com/crowdflux/angel/utilities"
	"github.com/crowdflux/angel/utilities"
)

//
//func TestHttpHit(t *testing.T) {
//	flus := []models.FeedLineUnit{models.FeedLineUnit{
//		ReferenceId: "PAYTM_123",
//		Data: models.JsonFake{
//			"product_id":  "da17b335-7ed9-4928-a222-44eda29a4896",
//			"category_id": "t_shirt_12",
//			"name":        "XYZ Men's Gold T-Shirt",
//			"brand":       "XYZ",
//			"color":       "Gold",
//		},
//		Tag: "PAYTM_TSHIRT",
//	}}
//
//	uuid, _ := uuid.FromString("da17b335-7ed9-4928-a222-44eda29a4896")
//	resp, status := sendBackToClient(uuid, flus)
//
//	fmt.Println("Status", status)
//	fmt.Println("Response", resp)
//	/*assert.NoError(t, status, "Error occured while validating")
//	assert.True(t, isValid, "Expected valid flu but found inValid")
//	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
//}
//
//func TestCallBack(t *testing.T) {
//	flus := []models.FeedLineUnit{models.FeedLineUnit{
//		ID:          uuid.NewV4(),
//		ReferenceId: "PAYTM_123",
//		Data: models.JsonFake{
//			"product_id":  "da17b335-7ed9-4928-a222-44eda29a4896",
//			"category_id": "t_shirt_12",
//			"name":        "XYZ Men's Gold T-Shirt",
//			"brand":       "XYZ",
//			"color":       "Gold",
//		},
//		Tag: "PAYTM_TSHIRT",
//	}}
//	id, _ := uuid.FromString("da17b335-7ed9-4928-a222-44eda29a4896")
//	feedLinePipe[id] = feedLineValue{utilities.TimeInMillis(), flus}
//	sendBackResp([]uuid.UUID{id})
//
//	//fmt.Println("Status",status)
//	/*assert.NoError(t, status, "Error occured while validating")
//	assert.True(t, isValid, "Expected valid flu but found inValid")
//	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
//}
//
//func TestBufferPut(t *testing.T) {
//	flus := []models.FeedLineUnit{models.FeedLineUnit{
//		ReferenceId: "PAYTM_123",
//		Data: models.JsonFake{
//			"product_id":  "59955f54-e75c-40a1-8d11-162e12dbf68a",
//			"category_id": "t_shirt_12",
//			"name":        "XYZ Men's Gold T-Shirt",
//			"brand":       "XYZ",
//			"color":       "Gold",
//		},
//		Tag: "PAYTM_TSHIRT",
//	}}
//
//	flm := &FluMonitor{}
//	resp := flm.AddManyToOutputQueue(flus)
//
//	//fmt.Println("Status",status)
//	fmt.Println("Response", resp)
//	/*assert.NoError(t, status, "Error occured while validating")
//	assert.True(t, isValid, "Expected valid flu but found inValid")
//	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
//}
//
//func TestStartFluOutputTimer(t *testing.T) {
//	StartFluOutputTimer()
//}
//
//func TestPutDbLog(t *testing.T) {
//	flus := []models.FeedLineUnit{models.FeedLineUnit{
//		ID:          uuid.NewV4(),
//		ReferenceId: "PAYTM_123",
//		Data: models.JsonFake{
//			"product_id":  "59955f54-e75c-40a1-8d11-162e12dbf68a",
//			"category_id": "t_shirt_12",
//			"name":        "XYZ Men's Gold T-Shirt",
//			"brand":       "XYZ",
//			"color":       "Gold",
//		},
//		Tag: "PAYTM_TSHIRT",
//	}}
//
//	flp := feedLineValue{utilities.TimeInMillis(), []models.FeedLineUnit{}}
//	for _, flu := range flus {
//		flp.feedLine = append(flp.feedLine, flu)
//	}
//	putDbLog(flp, "Success", Response{})
//}

func TestGetFluOutputObj(t *testing.T) {

	var i interface{} = models.JsonF{
		"1": 2,
	}

	flus := []models.FeedLineUnit{
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonF{
				RESULT: i,
			},
		},
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonF{
				RESULT: i,
			},
		},
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonF{
				RESULT: i,
			},
		},
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonF{
				RESULT: i,
			},
		},
	}
	flp := feedLineValue{10, utilities.TimeInMillis(), flus}

	fmt.Println(len(flp.feedLine))
	out := getFluOutputObj(flp)
	fmt.Println(len(out), out)

	projId := uuid.NewV4()
	feedLinePipe[projId] = flp

	fmt.Println(len(feedLinePipe[projId].feedLine))
	deleteFromFeedLinePipe(projId, out)
	fmt.Println(len(feedLinePipe[projId].feedLine))

}
