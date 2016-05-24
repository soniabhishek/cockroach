package config_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/config"
	"testing"
)

func Test(t *testing.T) {

	s := config.Get(config.PG_HOST)
	assert.EqualValues(t, "localhost", s)
}

func ExampleGet() {

	baseApi := config.Get(config.BASE_API_URL)

	fmt.Println(baseApi)
	// Output: localhost:8999/api
}

func ExampleIsDevelopment() {

	// Returns true if current env is development
	fmt.Println(config.IsDevelopment())
	// Output: true
}
