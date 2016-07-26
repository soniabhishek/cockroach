package auther_test

import (
	"fmt"

	"github.com/crowdflux/angel/app/api/auther"
	"github.com/crowdflux/angel/app/models/uuid"
)

func Example() {
	id := uuid.NewV4()
	key := auther.StdProdAuther.GetAPIKey(id)

	fmt.Println(id)
	fmt.Println(key)
}
