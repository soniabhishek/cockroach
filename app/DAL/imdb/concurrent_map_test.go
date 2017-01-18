package imdb

import (
	"fmt"
	"testing"
)

func TestNewFluValidateCache(t *testing.T) {
	c := cmap{}
	c.cmap = make(map[interface{}]interface{})
	c.set(1, 1)
	fmt.Println(c.get(1))
}

func Test(t *testing.T) {
	c := NewCmap()
	c.set(1, 1)
	c.set(2, 2)
	c.reset()

	fmt.Println(c.get(1))
	fmt.Println(c.get(2))

}
