package clients

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbacusSuccesfulPrediction(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err, success := cc.Predict("Good item")
	expectedResult := "Approve"
	assert.Equal(t, actualResult, expectedResult)
	assert.True(t, success, true)
	assert.NoError(t, err)
}

func TestAbacusUnsuccessfulPrediciton(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err, success := cc.Predict("Jango Fett")
	assert.NotEqual(t, actualResult, "Approve")
	assert.False(t, success)
	assert.NoError(t, err)
}
