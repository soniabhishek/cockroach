package clients

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbacusGood(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err, success := cc.Predict("Good item")
	expectedResult := "Approve"
	assert.Equal(t, actualResult, expectedResult)
	assert.True(t, success, true)
	assert.NoError(t, err)
}

func TestAbacusBad(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err, success := cc.Predict("Jango Fett")
	assert.NotEqual(t, actualResult, "Approve")
	assert.False(t, success)
	assert.NoError(t, err)
}
