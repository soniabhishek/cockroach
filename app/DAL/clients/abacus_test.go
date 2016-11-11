package clients

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbacus(t *testing.T) {
	cc := GetAbacusClient()
	actualResult, err := cc.Predict("Good item")
	expectedResult := "Approve"
	assert.Equal(t, actualResult, expectedResult)
	assert.NoError(t, err)
}
