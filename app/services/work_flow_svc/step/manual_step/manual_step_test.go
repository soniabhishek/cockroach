package manual_step

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models/uuid"
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
