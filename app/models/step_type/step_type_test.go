package step_type

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepType_Scan(t *testing.T) {
	s := CrowdSourcing
	val, err := s.Value()
	assert.NoError(t, err)
	assert.Equal(t, interface{}(uint(1)), val)

	err = s.Scan(int(1))
	assert.NoError(t, err)
	assert.Equal(t, CrowdSourcing, s)
}
