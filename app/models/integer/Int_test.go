package integer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type intTestStruct struct {
	in         int
	str        string
	shouldPass bool
}

var casesString = []intTestStruct{
	intTestStruct{2, "2", true},
	intTestStruct{200, "200", true},
	intTestStruct{300, "200", false},
	intTestStruct{100000, "100000", true},
	intTestStruct{100000000000000, "100000000000000", true},
}

var casesFromString = []intTestStruct{
	intTestStruct{2, "2", true},
	intTestStruct{200, "200", true},
	intTestStruct{300, "200", false},
	intTestStruct{100000, "100000", true},
	intTestStruct{100000000000000, "100000000000000", true},
}

func TestInt_String(t *testing.T) {
	i := Int(2)
	s := i.String()
	assert.Equal(t, "2", s)
}

func TestInt_FromString(t *testing.T) {
	i, err := FromString("2")
	assert.NoError(t, err)
	assert.Equal(t, Int(2), i)
	assert.Equal(t, 2, i.I())

	i, err = FromString("2o")
	assert.Error(t, err)

}

func TestInt_FromStringOrNil(t *testing.T) {
	i := FromStringOrNil("2")
	assert.Equal(t, Int(2), i)
	assert.Equal(t, 2, i.I())

	i = FromStringOrNil("2o")
	assert.Equal(t, Nil, i)
}

func BenchmarkFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromString("20")
	}
}

func BenchmarkFromStringOrNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromStringOrNil("20")
	}
}
