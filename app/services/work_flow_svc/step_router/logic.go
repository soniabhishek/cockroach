package step_router

import (
	"errors"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
)

var ErrLogicNotFound = errors.New("Logic not found")
var ErrLogicKeyNotFound = errors.New("logic key not found")
var ErrLogicKeyNotValid = errors.New("logic key not valid")

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

	// Future stuff below
	case "string_equality":
		//s := l.InputTemplate["value"].(StringEquality)
		//
		//s.Left = flu.Build[s.Left].(string)
		//s.Right = flu.Build[s.Right].(string)
		//return s.True(), nil
		fallthrough
	case "number_equality":
		//n := l.InputTemplate["value"].(NumberEquality)
		//n.Left = flu.Build[n.Left].(int)
		//n.Right = flu.Build[n.Right].(int)
		//return n.True(), nil
		fallthrough
	case "number_comparision":
		//n := l.InputTemplate["value"].(NumberComparison)
		//n.Left = flu.Build[n.Left].(int)
		//n.Right = flu.Build[n.Right].(int)
		//return n.True(), nil
		fallthrough
	default:
		return false, ErrLogicNotFound
	}
}

//--------------------------------------------------------------------------------//

type StringEquality struct {
	Left        string
	Right       string
	ShouldEqual bool
}

func (s *StringEquality) True() bool {
	if s.ShouldEqual {
		return s.Left == s.Right
	} else {
		return s.Left != s.Right
	}
}

type NumberEquality struct {
	Left        int
	Right       int
	ShouldEqual bool
}

func (n *NumberEquality) True() bool {
	if n.ShouldEqual {
		return n.Left == n.Right
	} else {
		return n.Left != n.Right
	}
}

type NumberComparison struct {
	Left          int
	Right         int
	ShouldGreater bool
}

func (n *NumberComparison) True() bool {
	if n.ShouldGreater {
		return n.Left > n.Right
	} else {
		return n.Left <= n.Right
	}
}
