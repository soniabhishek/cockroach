package auther

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func Test(t *testing.T) {
	id := uuid.FromStringOrNil("d1244145-bf59-4a05-8abd-254b41c0cafa")

	auth, err := New("playmentlnsdjvds")
	assert.Error(t, err)

	auth, err = New("playmentlnsdjvds som e uefghv khgvhh")
	assert.NoError(t, err)

	key := auth.GetAPIKey(id)

	tr := auth.Check(id, key)

	assert.True(t, tr)
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

func TestProdAuther(t *testing.T) {
	id := uuid.NewV4()
	s := StdProdAuther.GetAPIKey(id)
	fmt.Println(id.CEnc())
	fmt.Println(s)
}
