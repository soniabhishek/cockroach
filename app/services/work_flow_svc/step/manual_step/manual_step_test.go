package manual_step

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"os"
	"testing"
	"time"
)

func TestManualStep_DownloadCsv(t *testing.T) {
	manualStepId, _ := uuid.FromString("7adbafe4-1dc2-4ba3-b8ba-155c87666323")
	resp, err := DownloadCsv(manualStepId)

	fmt.Println("Err", err)
	fmt.Println("Response", resp)

	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}

func TestManualStep_UploadCsv(t *testing.T) {
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
			Build: models.JsonFake{
				"first": 123,
				//"second": models.JsonFake{
				//	"review": 123,
				//},
			},
		},
		models.FeedLineUnit{
			ID: uuid.NewV4(),
			Build: models.JsonFake{
				"asf": 432,
				//"second": models.JsonFake{
				//	"review": 1234,
				//},
			},
		},
	}

	fileString, _, err := createJSONFile(flus, "./", uuid.FromStringOrNil("f90f4e0c-c616-43ca-a83c-7d7b8dcf5bd5"))
	assert.NoError(t, err)

	plog.Info("create json test", fileString)

	//assert.Equal(t, )
}

func TestWriteFile(t *testing.T) {
	manualStepId := uuid.NewV4()
	file := TEMP_FOLDER + string(os.PathSeparator) + manualStepId.String() + ".txt"
	fmt.Println(manualStepId)
	err := createFile(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	csvBuff := megatronJson{make([]models.JsonFake, 0)}
	csvBuff.Jsons = append(csvBuff.Jsons, models.JsonFake{"One": 1})
	csvBuff.Jsons = append(csvBuff.Jsons, models.JsonFake{"Two": "2"})
	err = writeFile(file, csvBuff)
	if err != nil {
		fmt.Println(err.Error())
	}
}
