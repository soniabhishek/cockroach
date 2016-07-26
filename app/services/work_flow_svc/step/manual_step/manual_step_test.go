package manual_step

import (
	"testing"
	"time"

	"fmt"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

func TestManualStep_DownloadCsv(t *testing.T) {
	manualStepId, _ := uuid.FromString("d364cc9e-f267-4fa7-aca0-864f79f4a3c4")
	resp, err := DownloadCsv(manualStepId)

	assert.NoError(t, err)

	//fmt.Println("Err", err)
	fmt.Println("Response", resp)

	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}

func TestManualStep_UploadCsv(t *testing.T) {

	t.Skip("Api has been changed")

	err := UploadCsv("filename")
	assert.NoError(t, err)

	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}

func TestTime(t *testing.T) {

	timeStr := "2016-06-05T17:17:40+05:30"

	tym, err := time.Parse(timeFormat, timeStr)
	assert.NoError(t, err)
	assert.Equal(t, 17, tym.Hour())
	assert.Equal(t, 17, tym.Minute())
	assert.Equal(t, 40, tym.Second())
}

func TestCreateJsonFile(t *testing.T) {

	flus := []models.FeedLineUnit{
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonF{
				"first": 123,
				"second": models.JsonF{
					"review": 123,
				},
			},
		},
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonF{
				"asf": 432,
				"second": models.JsonF{
					"review": 1234,
				},
			},
		},
	}

	_, _, err := createJSONFile(flus, "./", uuid.FromStringOrNil("f90f4e0c-c616-43ca-a83c-7d7b8dcf5bd5"))
	assert.NoError(t, err)
}
