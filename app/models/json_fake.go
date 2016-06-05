package models

import (
	"database/sql/driver"
	"encoding/json"
)

//Need to reimplement this type
//The value & scan methods needs to be properly written
type JsonFake map[string]interface{}

// Value returns a driver Value.
func (j JsonFake) Value() (driver.Value, error) {

	return j.String(), nil
}

// An error should be returned if the value can not be stored
// without loss of information.
func (j *JsonFake) Scan(src interface{}) error {

	var tmp map[string]interface{}

	bty, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bty, &tmp)
	if err != nil {
		return err
	}
	*j = JsonFake(tmp)
	return nil
}

/**
The is a helper function
*/
func (j *JsonFake) String() string {
	/**
	bool, for JSON booleans
	float64, for JSON numbers
	string, for JSON strings
	[]interface{}, for JSON arrays
	map[string]interface{}, for JSON objects
	nil for JSON null
	*/

	bty, _ := json.Marshal(*j)
	return string(bty)
}

func (j *JsonFake) Merge(a JsonFake) {
	for k, v := range a {
		(*j)[k] = v
	}
}
