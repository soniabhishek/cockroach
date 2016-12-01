package clients

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// WORKS ONLY 10% OF THE TIME DUE TO 10% SUCCESS LOGIC IN ABACUS
func TestAbacusSuccesfulPrediction(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err, success := cc.Predict("Good item")
	expectedResult := "Approve-Useful Product Review"
	assert.Equal(t, expectedResult, actualResult)
	assert.True(t, success)
	assert.NoError(t, err)
}

func TestAbacusUnsuccessfulPrediciton(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err, success := cc.Predict("Jango Fett")
	assert.NotEqual(t, actualResult, "Approve")
	assert.False(t, success)
	assert.NoError(t, err)
}
