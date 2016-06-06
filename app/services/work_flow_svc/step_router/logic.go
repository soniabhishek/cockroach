package step_router

import (
	"errors"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

var ErrLogicNotFound = errors.New("Logic not found")
var ErrLogicTypeNotFound = errors.New("Type not found")
var ErrLogicTypeNotValid = errors.New("Type not valid")

func Logic(flu feed_line.FLU, l models.LogicGate) (bool, error) {

	templateType, ok := l.InputTemplate["type"]
	if !ok {
		return false, ErrLogicTypeNotFound
	}

	templateTypeStr, ok := templateType.(string)
	if !ok {
		return false, ErrLogicTypeNotValid
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
