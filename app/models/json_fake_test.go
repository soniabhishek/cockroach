package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonFake_String(t *testing.T) {
	var jsn = JsonF{
		"fir": "st",
		"sec": true,
	}
	bty, err := json.Marshal(jsn)
	assert.NoError(t, err)

	jsnStr := string(bty)
	assert.Equal(t, jsnStr, jsn.String())
}

func TestJsonFake_Scan(t *testing.T) {

	jsn := JsonF{
		"fir": "st",
		"sec": "ond",
		"third": JsonF{
			"inner": "peace",
		},
	}

	mapStrInterf := map[string]interface{}{
		"fir": "st",
		"sec": "ond",
		"third": map[string]interface{}{
			"inner": "peace",
		},
	}

	jsn2 := JsonF{}
	mapStrBty, err := json.Marshal(mapStrInterf)

	err = jsn2.Scan(mapStrBty)
	assert.NoError(t, err)
	assert.EqualValues(t, jsn.String(), jsn2.String())
}

func TestJsonFake_Merge(t *testing.T) {
	jsn := JsonF{
		"a": 1,
	}
	jsn2 := JsonF{
		"b": 2,
	}
	jsn.Merge(jsn2)
	assert.EqualValues(t, jsn, JsonF{
		"a": 1,
		"b": 2,
	})
}
func TestJsonF_Scan(t *testing.T) {
	jsn := map[string]interface{}{
		"a": 1,
	}

	j1 := JsonF(jsn)
	j2 := JsonF{}
	err := j2.Scan(jsn)
	assert.NoError(t, err)
	assert.Equal(t, j1, j2)
}
