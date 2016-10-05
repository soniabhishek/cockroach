package router_svc

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/stretchr/testify/assert"
)

type logicGateTestCase struct {
	LogicGate models.LogicGate
	Result    bool
	Error     error
}

func TestLogic_Continue(t *testing.T) {

	var logicGateTestCases = []logicGateTestCase{
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

func TestLogic_Boolean(t *testing.T) {

	var logicGateTestCases = []logicGateTestCase{
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "boolean",
				},
			},
			Result: false,
			Error:  ErrMalformedLogicOptions,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "boolean",
					"options": map[string]interface{}{
						"blah": "blah",
					},
				},
			},
			Result: false,
			Error:  ErrMalformedLogicOptions,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "boolean",
					"options": map[string]interface{}{
						"field_name":     "blah",
						"should_be_true": "asd",
					},
				},
			},
			Result: false,
			Error:  ErrMalformedLogicOptions,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "boolean",
					"options": map[string]interface{}{
						"field_name":     "blah",
						"should_be_true": true,
					},
				},
			},
			Result: false,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "boolean",
					"options": map[string]interface{}{
						"field_name":     "abcd",
						"should_be_true": true,
					},
				},
			},
			Result: false,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "boolean",
					"options": map[string]interface{}{
						"field_name":     "pqrs",
						"should_be_true": true,
					},
				},
			},
			Result: true,
			Error:  nil,
		},
	}

	flu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			Build: models.JsonF{
				"abcd": 1,
				"pqrs": true,
			},
		},
	}
	for i, testCase := range logicGateTestCases {

		out, err := Logic(flu, testCase.LogicGate)
		assert.Equal(t, testCase.Error, err, "index:", i)
		assert.EqualValues(t, testCase.Result, out, "index:", i)
	}

}

func TestLogic_StringEquality(t *testing.T) {

	var logicGateTestCases = []logicGateTestCase{
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "string_equal",
				},
			},
			Result: false,
			Error:  ErrMalformedLogicOptions,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "string_equal",
					"options": map[string]interface{}{
						"field_name":      "mnop",
						"should_be_equal": true,
						"field_value":     "Hello",
					},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "string_equal",
					"options": map[string]interface{}{
						"field_name":      "mnop",
						"should_be_equal": false,
						"field_value":     "Lollo",
					},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"logic": "string_equal",
					"options": map[string]interface{}{
						"field_name":      "mnop",
						"should_be_equal": false,
						"field_value":     map[string]interface{}{"asd": 1},
					},
				},
			},
			Result: false,
			Error:  ErrMalformedLogicOptions,
		},
	}

	flu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			Build: models.JsonF{
				"abcd": 1,
				"pqrs": true,
				"mnop": "Hello",
			},
		},
	}
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
