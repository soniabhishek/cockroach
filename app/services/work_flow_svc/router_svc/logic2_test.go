package router_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogic2(t *testing.T) {

	var logicGateTestCases = []logicGateTestCase{
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"options": map[string]interface{}{
						"logic_expression": "{abcd} >1 || {xyz} == 'man' || IsNull({arc})",
					},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"options": map[string]interface{}{
						"logic_expression": "{pqrs} && {abcd} <3 && {xyz} == 'man' && !IsNull({arc})",
					},
				},
			},
			Result: false,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: models.LogicGate{
				InputTemplate: models.JsonF{
					"options": map[string]interface{}{
						"logic_expression": "{pqrs} && {abCD} <3 && {xyz} == 'man' && !IsNull({arc})",
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
				"pqrs": false,
				"xyz":  "man",
			},
		},
	}
	for i, testCase := range logicGateTestCases {

		out, err := Logic2(flu, testCase.LogicGate)
		assert.Equal(t, testCase.Error, err, "index:", i)
		assert.EqualValues(t, testCase.Result, out, "index:", i)
	}

}

func BenchmarkLogic2(b *testing.B) {

	logicGate := models.LogicGate{
		InputTemplate: models.JsonF{
			"options": map[string]interface{}{
				"logic_expression": "abcd >1 || xyz == 'man' || IsNull(arc)",
				"logic_fields":     []string{"abcd", "xyz", "arc"},
			},
		},
	}

	for i := 0; i < b.N; i++ {

		flu := feed_line.FLU{
			FeedLineUnit: models.FeedLineUnit{
				Build: models.JsonF{
					"abcd": 1,
					"pqrs": false,
					"xyz":  "man",
					"arc":  1,
				},
			},
		}
		Logic2(flu, logicGate)
	}

}
