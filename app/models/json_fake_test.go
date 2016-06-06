package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonFake_String(t *testing.T) {
	var jsn = JsonFake{
		"fir": "st",
		"sec": true,
	}
	bty, err := json.Marshal(jsn)
	assert.NoError(t, err)

	jsnStr := string(bty)
	assert.Equal(t, jsnStr, jsn.String())
}

func TestJsonFake_Scan(t *testing.T) {

	jsn := JsonFake{
		"fir": "st",
		"sec": "ond",
		"third": JsonFake{
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

	jsn2 := JsonFake{}
	mapStrBty, err := json.Marshal(mapStrInterf)

	err = jsn2.Scan(mapStrBty)
	assert.NoError(t, err)
	assert.EqualValues(t, jsn.String(), jsn2.String())
}

func TestJsonFake_Merge(t *testing.T) {
	jsn := JsonFake{
		"a": 1,
	}
	jsn2 := JsonFake{
		"b": 2,
	}
	jsn.Merge(jsn2)
	assert.EqualValues(t, jsn, JsonFake{
		"a": 1,
		"b": 2,
	})
}
