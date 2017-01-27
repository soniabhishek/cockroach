package imdb

import (
	"fmt"
	"testing"
)

func TestNewFluValidateCache(t *testing.T) {
	c := cmap{}
	c.cmap = make(map[interface{}]interface{})
	c.Set(1, 1)
	fmt.Println(c.Get(1))
}

func Test(t *testing.T) {
	c := NewCmap()
	c.Set(1, 1)
	c.Set(2, 2)
	c.Reset()

	fmt.Println(c.Get(1))
	fmt.Println(c.Get(2))

}

func TestCmap_Iter(t *testing.T) {
	c := cmap{}
	c.cmap = make(map[interface{}]interface{})
	c.Set("abhishek", "1")
	c.Set("soni", "2")
	c.Set("jaipur", "3")
	for tuple := range c.Iter() {
		fmt.Println(tuple)
	}
}
