package step_router

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"testing"
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
				InputTemplate: models.JsonFake{
					"type": "continue",
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonFake{
					"type": "continues",
				},
			},
			Result: false,
			Error:  ErrLogicNotFound,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonFake{
					"type12": "continue",
				},
			},
			Result: false,
			Error:  ErrLogicTypeNotFound,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonFake{
					"type": 123,
				},
			},
			Result: false,
			Error:  ErrLogicTypeNotValid,
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
