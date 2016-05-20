package integer

import (
	"strconv"
)

type Int int

var Nil Int

/**
returns int -> string
*/
func (i *Int) String() string {
	return strconv.Itoa(int(*i))
}

/**
Gets the underlying primitive int of the integer.Int
*/
func (i *Int) I() int {
	return int(*i)
}

func FromString(s string) (o Int, err error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	return Int(i), nil
}

func FromStringOrNil(s string) (o Int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return Nil
	}
	return Int(i)
}
