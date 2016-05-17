package models

import (
	"database/sql/driver"
	"strconv"
)

//Need to reimplement this type
//The value & scan methods needs to be properly written
type JsonFake map[string]interface{}

// Value returns a driver Value.
func (j JsonFake) Value() (driver.Value, error) {

	//fmt.Println("value called", j.String())

	return j.String(), nil
}

// An error should be returned if the value can not be stored
// without loss of information.
func (j *JsonFake) Scan(src interface{}) error {
	return nil
}

/**
The is a helper function
*/
func (j *JsonFake) String() string {

	x := "{"

	for key, value := range *j {

		switch value.(type) {
		case string:
			x += "\"" + key + "\" : \"" + value.(string) + "\","
		case float64:
			x += "\"" + key + "\" : \"" + strconv.FormatFloat(value.(float64), 'E', -1, 64) + "\","
		}
	}

	newVal := x

	//Remove the last ',' character
	if x != "{" {
		newVal = x[:len(x)-1]
	}

	//Add the enclosing bracket
	newVal += "}"

	return newVal
}
