package config_test

import (
	"fmt"
	"testing"

	"github.com/crowdflux/angel/app/config"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

	s := config.PG_HOST.Get()
	assert.EqualValues(t, "localhost", s)
}

func ExampleGet() {

	baseApi := config.BASE_API_URL.Get()

	fmt.Println(baseApi)
	// Output: localhost:8999/api
}

func ExampleIsDevelopment() {

	// Returns true if current env is development
	fmt.Println(config.IsDevelopment())
	// Output: true
}
