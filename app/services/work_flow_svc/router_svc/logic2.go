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

	fields := getFieldsFromExpression(expRegEx, exp)

	//trimming actual expression of curly braces
	exp = strings.Replace(strings.Replace(exp, "{", "", -1), "}", "", -1)

	parameters := getParametersFromFlu(flu, fields)

	expression, err := getEvaluatableExpression(exp, flu, fields)
	if err != nil {
		return false, err
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

// extracting variables from the expression using reg ex
func getFieldsFromExpression(regEx string, exp string) (fields []string) {
	re := regexp.MustCompile(regEx)
	fields = re.FindAllString(exp, -1)

	//trimming the expression and fields of curly braces
	for i, field := range fields {
		fields[i] = strings.TrimLeft(strings.TrimRight(field, "}"), "{")
	}
	return
}

func getParametersFromFlu(flu feed_line.FLU, fields []string) (parameters map[string]interface{}) {
	//getting the actual parameters from flu using keys
	parameters = make(map[string]interface{}, len(fields))
	for _, field := range fields {
		parameters[field] = flu.Build[field]
	}
	return
}

func getEvaluatableExpression(exp string, flu feed_line.FLU, fields []string) (expression *govaluate.EvaluableExpression, err error) {
	//get from cache if it's there
	expression, ok := expressionCache[exp]
	if !ok {
		expression, err = govaluate.NewEvaluableExpressionWithFunctions(exp, getCustomFunctions(flu, fields))
		if err != nil {
			return expression, ErrMalformedLogicOptions
		}
		expressionCache[exp] = expression
	}
	return expression, nil
}
