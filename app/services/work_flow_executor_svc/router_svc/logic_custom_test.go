package router_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogicCustom(t *testing.T) {

	var logicGateTestCases = []logicGateTestCase{
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "StringContains({xyz},ToUpper('ma'),'MAN')",
					}},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "StringContains({efgh},1)",
					}},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "StringLength({xyz})==3",
					}},
				},
			},
			Result: false,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "IsNull({arc})",
					}},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "IsNull({arc}) || {arc}==2",
					}},
				},
			},
			Result: false,
			Error:  ErrPropNotFoundInFluBuild,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "IsNull({arc})?false:{arc}>3",
					}},
				},
			},
			Result: false,
			Error:  ErrPropNotFoundInFluBuild,
		},

		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "IsNull({nostring})",
					}},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "IsNull({nilfield})",
					}},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "IsNull({efgh})",
					}},
				},
			},
			Result: false,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "{bcd}==1",
					}},
				},
			},
			Result: false,
			Error:  ErrPropNotFoundInFluBuild,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "(ToUpper({ijkl})+ToLower({xyz}))=='DOGman'",
					}},
				},
			},
			Result: true,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "{pqrs} && {abcd} <3 && {xyz} == 'MAN' && !IsNull({arc})",
					}},
				},
			},
			Result: false,
			Error:  nil,
		},
		logicGateTestCase{
			LogicGate: LogicGate{
				InputTemplate: []models.JsonF{{
					"logic": "custom",
					"options": map[string]interface{}{
						expression_field: "{pqrs} && {abCD} <3 && {xyz} in ('man','woman') && !IsNull({arc})",
					}},
				},
			},
			Result: false,
			Error:  ErrPropNotFoundInFluBuild,
		},
	}
	flu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			Build: models.JsonF{
				"abcd":     1,
				"efgh":     "GOD1",
				"pqrs":     false,
				"xyz":      "MAN",
				"ijkl":     "dog",
				"nostring": "",
				"nilfield": nil,
			},
		},
	}
	for i, testCase := range logicGateTestCases {

		out, err := EvaluateLogics(flu, testCase.LogicGate)
		assert.Equal(t, testCase.Error, err, "index:", i)
		assert.EqualValues(t, testCase.Result, out, "index:", i)
	}

}

func BenchmarkLogic2(b *testing.B) {

	logicGate := LogicGate{
		InputTemplate: []models.JsonF{{
			"logic": "custom",
			"options": map[string]interface{}{
				expression_field: "abcd >1 || xyz == 'man' || IsNull(arc)",
			}},
		},
	}

	for i := 0; i < b.N; i++ {

		flu := feed_line.FLU{
			FeedLineUnit: models.FeedLineUnit{
				Build: models.JsonF{
					"abcd": 1,
					"pqrs": false,
					"xyz":  "MAN",
					"arc":  1,
				},
			},
		}
		EvaluateLogics(flu, logicGate)
	}

}
