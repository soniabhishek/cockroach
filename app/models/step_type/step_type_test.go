package step_type

import (
	"testing"

	"github.com/notnow/src/encoding/json"
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

func TestStepType_Value(t *testing.T) {

	type Name struct {
		A StepType `json:"A"`
	}

	str := `{ "A" : 1 }`

	n := Name{}

	err := json.Unmarshal([]byte(str), &n)
	assert.NoError(t, err)
}
