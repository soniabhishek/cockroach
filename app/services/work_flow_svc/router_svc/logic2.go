package router_svc

import (
	"github.com/Knetic/govaluate"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"regexp"
	"strings"
)

// to cache the expressions to be reused
var expressionCache = make(map[string]*govaluate.EvaluableExpression)

func Logic2(flu feed_line.FLU, l models.LogicGate) (value interface{}, err error) {

	options, ok1 := l.InputTemplate["options"].(map[string]interface{})
	exp, ok2 := options["logic_expression"].(string)
	expRegEx := "{[A-z,0-9,_,-]*}"
	if !ok1 || !ok2 {
		return false, ErrMalformedLogicOptions
	}

	// extracting variables from the expression using reg ex
	re := regexp.MustCompile(expRegEx)
	fields := re.FindAllString(exp, -1)

	//trimming the expression and fields of curly braces
	for i, field := range fields {
		fields[i] = strings.TrimLeft(strings.TrimRight(field, "}"), "{")
	}
	exp = strings.Replace(strings.Replace(exp, "{", "", -1), "}", "", -1)

	//getting the actual parameters from flu using keys
	parameters := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		parameters[field] = flu.Build[field]
	}

	expression, ok := expressionCache[exp]
	if !ok {
		expression, err = govaluate.NewEvaluableExpressionWithFunctions(exp, getCustomFunctions(flu, fields))
		if err != nil {
			return false, ErrMalformedLogicOptions
		}
		expressionCache[exp] = expression
	}

	value, err = expression.Evaluate(parameters)
	if err != nil {
		return false, ErrMalformedLogicOptions
	}
	return
}

// more custom functions can be added in future
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
