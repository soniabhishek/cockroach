package step_router

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
	"github.com/stretchr/testify/assert"
)

func TestLogic(t *testing.T) {

	type logicGateTestCase struct {
		LogicGate models.LogicGate
		Result    bool
		Error     error
	}
	logicGateTestCases := []logicGateTestCase{
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "continue",
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "continues",
				},
			},
			Result: false,
			Error:  ErrLogicNotFound,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic12": "continue",
				},
			},
			Result: false,
			Error:  ErrLogicKeyNotFound,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": 123,
				},
			},
			Result: false,
			Error:  ErrLogicKeyNotValid,
		},
	}

	var flu feed_line.FLU
	for i, testCase := range logicGateTestCases {

		out, err := Logic(flu, testCase.LogicGate)
		assert.Equal(t, testCase.Error, err, "index:", i)
		assert.EqualValues(t, testCase.Result, out, "index:", i)
	}

}

func TestRandom(t *testing.T) {
	type sSs struct {
		Left        string
		Right       string
		ShouldEqual bool
	}
	jsn := `{
		"left" : "haramkhor",
		"right" : "chuitya",
		"shouldequal: true
	}`

	var inter interface{} = jsn

	ss, ok := inter.(sSs)
	fmt.Println(ss)
	fmt.Println(ok)
}
