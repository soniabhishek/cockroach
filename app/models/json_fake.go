package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//Need to reimplement this type
//The value & scan methods needs to be properly written
type JsonF map[string]interface{}

// Value returns a driver Value.
func (j JsonF) Value() (driver.Value, error) {

	return j.String(), nil
}

// An error should be returned if the value can not be stored
// without loss of information.
func (j *JsonF) Scan(src interface{}) error {

	var tmp map[string]interface{}
	var bty []byte

	switch src.(type) {
	case []byte:
		bty = src.([]byte)
	case string:
		bty = []byte(src.(string))
	case map[string]interface{}:
		*j = JsonF(src.(map[string]interface{}))
		return nil
	default:
		return errors.New("only []byte & string supported at the moment")
	}

	err := json.Unmarshal(bty, &tmp)
	if err != nil {
		return err
	}
	*j = JsonF(tmp)
	return nil
}

/**
The is a helper function
*/
func (j *JsonF) String() string {
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

func (j *JsonF) Merge(a JsonF) JsonF {

	if *j == nil {
		*j = JsonF{}
	}

	for k, v := range a {
		(*j)[k] = v
	}
	return (*j)
}

func (j *JsonF) StringPretty() string {
	bty, _ := json.MarshalIndent(*j, "", "  ")
	return string(bty)
}
