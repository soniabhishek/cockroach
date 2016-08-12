package step_router

import (
	"errors"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
)

var ErrLogicNotFound = errors.New("Logic not found")
var ErrLogicKeyNotFound = errors.New("logic key not found")
var ErrLogicKeyNotValid = errors.New("logic key not valid")
var ErrMalformedLogicOptions = errors.New("Malformed logic options")

//var ErrPropNotFoundInFluBuild = errors.New("property not found in flu build")

func Logic(flu feed_line.FLU, l models.LogicGate) (bool, error) {

	templateType, ok := l.InputTemplate["logic"]
	if !ok {
		return false, ErrLogicKeyNotFound
	}

	templateTypeStr, ok := templateType.(string)
	if !ok {
		return false, ErrLogicKeyNotValid
	}

	switch templateTypeStr {
	case "continue":
		return true, nil
	case "boolean":

		options, ok1 := l.InputTemplate["options"].(map[string]interface{})
		shouldBeTrue, ok2 := options["should_be_true"].(bool)
		fieldName, ok3 := options["field_name"].(string)

		if !ok1 || !ok2 || !ok3 {
			return false, ErrMalformedLogicOptions
		}

		// ignoring field not present case in flu build
		// as it will return zero value of boolean which is false
		fieldValue, ok := flu.Build[fieldName].(bool)
		if !ok {
			plog.Trace("logic gate", "field not found for fluid ", flu.ID)
		}

		b := Boolean{
			FieldValue:   fieldValue,
			ShouldBeTrue: shouldBeTrue,
		}
		return b.True(), nil
	default:
		return false, ErrLogicNotFound
	}

}

//--------------------------------------------------------------------------------//

type Boolean struct {
	FieldValue   bool
	ShouldBeTrue bool
}

func (b *Boolean) True() bool {
	if b.ShouldBeTrue {
		return b.FieldValue == true
	} else {
		return b.FieldValue == false
	}
}

//==================================================================

type StringEquality struct {
	Left          string
	Right         string
	ShouldBeEqual bool
}

func (s *StringEquality) True() bool {
	if s.ShouldBeEqual {
		return s.Left == s.Right
	} else {
		return s.Left != s.Right
	}
}

type NumberEquality struct {
	Left          int
	Right         int
	ShouldBeEqual bool
}

func (n *NumberEquality) True() bool {
	if n.ShouldBeEqual {
		return n.Left == n.Right
	} else {
		return n.Left != n.Right
	}
}

type NumberComparison struct {
	Left            int
	Right           int
	ShouldBeGreater bool
}

func (n *NumberComparison) True() bool {
	if n.ShouldBeGreater {
		return n.Left > n.Right
	} else {
		return n.Left <= n.Right
	}
}
