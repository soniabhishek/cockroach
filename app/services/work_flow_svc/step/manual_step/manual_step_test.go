package manual_step

import (
	"fmt"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"testing"
)

func TestManualStep_DownloadCsv(t *testing.T) {
	ms := manualStep{}
	uuid1, _ := uuid.FromString("59955f54-e75c-40a1-8d11-162e12dbf68c")
	uuid2, _ := uuid.FromString("59955f54-e75c-40a1-8d11-162e12dbf68z")
	resp := ms.DownloadCsv(uuid1, uuid2)

	fmt.Println("Response", resp)
	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}

func TestManualStep_UploadCsv(t *testing.T) {
	ms := manualStep{}
	uuid1, _ := uuid.FromString("59955f54-e75c-40a1-8d11-162e12dbf68c")
	uuid2, _ := uuid.FromString("59955f54-e75c-40a1-8d11-162e12dbf68z")
	ms.UploadCsv("./", uuid1, uuid2)

	/*assert.NoError(t, status, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")*/
}
