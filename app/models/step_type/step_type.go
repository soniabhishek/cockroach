package step_type

import (
	"database/sql/driver"
	"errors"
)

// Kinda step enums
type StepType uint

const (
	CrowdSourcing StepType = iota + 1
	InternalSourcing
	Transformation
	Algorithm
	Bifurcation
	Unification
	Manual
	Gateway
	Error
	StartStep
	Test
)

func (s *StepType) Value() (driver.Value, error) {
	return uint(*s), nil
}
func (s *StepType) Scan(src interface{}) error {

	var tmp uint

	switch src.(type) {
	case uint:
		tmp = src.(uint)
	case int:
		tmp = uint(src.(int))
	case int64:
		tmp = uint(src.(int64))
	default:
		return errors.New("Not supported type passed")
	}
	*s = StepType(tmp)
	return nil
}

var stepTypeNames = map[StepType]string{
	CrowdSourcing:    "CrowdSourcing",
	InternalSourcing: "InternalSourcing",
	Transformation:   "Transformation",
	Algorithm:        "Algorithm",
	Bifurcation:      "Bifurcation",
	Unification:      "Unification",
	Manual:           "Manual",
	Gateway:          "Gateway",
	Error:            "Error",
	StartStep:        "StartStep",
	Test:             "Test",
}

func (s *StepType) String() string {
	if name, ok := stepTypeNames[*s]; ok {
		return name
	}
	return "NoName"
}
