package router_svc

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
)

var expressionCache = make(map[string]*govaluate.EvaluableExpression)

func Logic2(flu feed_line.FLU, l models.LogicGate) (value interface{}, err error) {

	options, ok1 := l.InputTemplate["options"].(map[string]interface{})
	exp, ok2 := options["logic_expression"].(string)
	fields, ok3 := options["logic_fields"].([]string)

	if !ok1 || !ok2 || !ok3 {
		fmt.Println("232324", ok1, ok2, ok3)
		return false, ErrMalformedLogicOptions
	}

	parameters := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		parameters[field] = flu.Build[field]
	}

	expression, ok := expressionCache[exp]

	if !ok {
		expression, err = govaluate.NewEvaluableExpressionWithFunctions(exp, getCustomFunctions(flu, fields))
		if err != nil {
			return
		}
		expressionCache[exp] = expression
	}

	value, err = expression.Evaluate(parameters)
	if err != nil {
		return false, ErrMalformedLogicOptions
	}
	return
}

func getCustomFunctions(flu feed_line.FLU, fields []string) map[string]govaluate.ExpressionFunction {
	funcs := map[string]govaluate.ExpressionFunction{
		"IsNull": func(params ...interface{}) (interface{}, error) {
			for _, param := range fields {
				if flu.Build[param] == nil {
					return true, nil
				}
			}
			return false, nil
		},
	}
	return funcs
}
