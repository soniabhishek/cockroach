package router_svc

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"regexp"
	"strings"
)

// to cache the expressions to be reused
var expressionCache = make(map[string]*govaluate.EvaluableExpression)

// to cache the fields to be reused
var fieldCache = make(map[string][]string)

var expression_field = "expression"

var expRegEx = "{((.*?))}"

func LogicCustom(flu feed_line.FLU, l models.LogicGate) (value bool, err error) {

	options, ok1 := l.InputTemplate["options"].(map[string]interface{})
	exp, ok2 := options[expression_field].(string)
	if !ok1 || !ok2 {
		return false, ErrMalformedLogicOptions
	}

	fields := getFieldsFromExpression(exp)

	//trimming actual expression of curly braces
	exp = strings.Replace(strings.Replace(exp, "{", "", -1), "}", "", -1)

	parameters, err := getParametersFromFlu(flu, fields, exp)
	if err != nil {
		return false, err
	}

	expression, err := getEvaluatableExpression(exp)
	if err != nil {
		return false, err
	}

	val, err := expression.Evaluate(parameters)

	if err != nil || val == nil {
		return false, ErrMalformedLogicOptions
	}

	value = val.(bool)
	return
}

// more custom functions can be added in future
func getCustomFunctions() map[string]govaluate.ExpressionFunction {
	funcs := map[string]govaluate.ExpressionFunction{
		"IsNull": func(params ...interface{}) (interface{}, error) {
			if len(params) > 0 {
				return false, nil
			}
			return true, nil
		},

		"ToLower": func(params ...interface{}) (interface{}, error) {
			if len(params) == 0 {
				return false, ErrPropNotFoundInFluBuild
			}

			if (len(params)) > 1 {
				return false, ErrMalformedLogicOptions
			}

			switch params[0].(type) {
			case string:
				return strings.ToLower(params[0].(string)), nil
			default:
				return params[0], nil
			}
		},
		"ToUpper": func(params ...interface{}) (interface{}, error) {
			if len(params) == 0 {
				return false, ErrPropNotFoundInFluBuild
			}

			if (len(params)) > 1 {
				return false, ErrMalformedLogicOptions
			}

			switch params[0].(type) {
			case string:
				return strings.ToUpper(params[0].(string)), nil
			default:
				return params[0], nil
			}
		},
		"StringContains": func(params ...interface{}) (interface{}, error) {
			if len(params) == 1 {
				return false, ErrMalformedLogicOptions
			}

			str := params[0].(string)
			strLen := len(params[0].(string))
			for _, param := range params {

				switch param.(type) {
				case string:
					p := param.(string)
					paramLen := len(p)
					if paramLen > strLen || !strings.Contains(str, p) {
						return false, nil
					}
				case float64:
					p := fmt.Sprint(param)
					paramLen := len(p)
					if paramLen > strLen || !strings.Contains(str, p) {
						return false, nil
					}
				default:
					return false, ErrMalformedLogicOptions
				}
			}
			return true, nil

		},
		"StringLength": func(params ...interface{}) (interface{}, error) {
			if len(params) == 0 {
				return false, ErrPropNotFoundInFluBuild
			}

			if len(params) > 1 {
				return false, ErrMalformedLogicOptions
			}
			switch params[0].(type) {
			case string:
				return len(params[0].(string)), nil
			default:
				return false, ErrMalformedLogicOptions
			}
		},
	}
	return funcs
}

// extracting variables from the expression using reg ex
func getFieldsFromExpression(exp string) (fields []string) {
	fields, ok := fieldCache[exp]
	if ok {
		return
	}
	re := regexp.MustCompile(expRegEx)
	fields = re.FindAllString(exp, -1)

	//trimming the fields of curly braces
	for i, field := range fields {
		fields[i] = strings.TrimLeft(strings.TrimRight(field, "}"), "{")
	}
	fieldCache[exp] = fields
	return
}

func getParametersFromFlu(flu feed_line.FLU, fields []string, exp string) (parameters map[string]interface{}, err error) {
	//getting the actual parameters from flu using keys
	parameters = make(map[string]interface{}, len(fields))
	for _, field := range fields {
		parameters[field] = flu.Build[field]
		if parameters[field] == nil && !strings.Contains(exp, "IsNull("+field) {
			return nil, ErrPropNotFoundInFluBuild
		}
	}
	return
}

func getEvaluatableExpression(exp string) (expression *govaluate.EvaluableExpression, err error) {
	//get from cache if it's there
	expression, ok := expressionCache[exp]
	if !ok {
		expression, err = govaluate.NewEvaluableExpressionWithFunctions(exp, getCustomFunctions())
		if err != nil {
			return expression, ErrMalformedLogicOptions
		}
		expressionCache[exp] = expression
	}
	return expression, nil
}
