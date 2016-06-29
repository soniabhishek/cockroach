package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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
	var bty []byte

	switch src.(type) {
	case []byte:
		bty = src.([]byte)
	case string:
		bty = []byte(src.(string))
	default:
		return errors.New("only []byte & string supported at the moment")
	}

	err := json.Unmarshal(bty, &tmp)
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

func (j *JsonFake) StringPretty() string {
	bty, err := json.MarshalIndent(*j, "", "  ")
	if err != nil {
		fmt.Println("JsonFakError:", err)
	}
	return string(bty)
}
