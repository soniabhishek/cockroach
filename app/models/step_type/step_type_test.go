package step_type

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
