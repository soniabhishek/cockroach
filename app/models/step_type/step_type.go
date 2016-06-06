package step_type

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
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
)

func (s StepType) Value() (driver.Value, error) {
	return uint(s), nil
}
func (s *StepType) Scan(src interface{}) error {

	var tmp uint

	fmt.Println("aasdsad", reflect.TypeOf(src).Kind().String())

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
