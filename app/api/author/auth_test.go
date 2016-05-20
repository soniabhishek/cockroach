package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"testing"
)

func Test(t *testing.T) {
	id := uuid.NewV4()

	auth, err := New("playmentlnsdjvds")
	assert.NoError(t, err)

	key := auth.GetAPIKey(id)

	tr := auth.Check(id, key)

	assert.True(t, tr)
	fmt.Println(id)
	fmt.Println(key)
	fmt.Println(tr)
}

var someId = uuid.NewV4()

func BenchmarkAuthor_GetAPIKey(b *testing.B) {

	for n := 0; n < b.N; n++ {

		id := someId
		auth, _ := New("playmentlnsdjvds")
		_ = auth.GetAPIKey(id)

	}
}
func BenchmarkAuthor_Check(b *testing.B) {

	for n := 0; n < b.N; n++ {
		id := someId
		auth, _ := New("playmentlnsdjvds")
		key := auth.GetAPIKey(id)
		auth.Check(id, key)
	}
}
