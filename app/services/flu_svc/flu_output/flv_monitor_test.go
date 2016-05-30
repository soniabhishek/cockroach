package flu_output

import (
	"fmt"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/utilities"
	"testing"
)

func TestHttpHit(t *testing.T) {
	flus := []models.FeedLineUnit{models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "59955f54-e75c-40a1-8d11-162e12dbf68a",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}}

	uuid, _ := uuid.FromString("59955f54-e75c-40a1-8d11-162e12dbf68a")
	resp, status := sendBackToClient(uuid, flus)

	fmt.Println("Status", status)
	fmt.Println("Response", resp)
	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}

func TestCallBack(t *testing.T) {
	flus := []models.FeedLineUnit{models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "59955f54-e75c-40a1-8d11-162e12dbf68a",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}}
	id, _ := uuid.FromString("59955f54-e75c-40a1-8d11-162e12dbf68a")
	feedLinePipe[id] = feedLineValue{utilities.TimeInMillis(), flus}
	sendBackResp([]uuid.UUID{id})

	//fmt.Println("Status",status)
	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}

func TestBufferPut(t *testing.T) {
	flus := []models.FeedLineUnit{models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "59955f54-e75c-40a1-8d11-162e12dbf68a",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}}

	flm := &FluMonitor{}
	resp := flm.AddManyToOutputQueue(flus)

	//fmt.Println("Status",status)
	fmt.Println("Response", resp)
	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}
