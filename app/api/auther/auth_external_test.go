package auther_test

import (
	"fmt"
	"gitlab.com/playment-main/angel/app/api/auther"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func Example() {
	id := uuid.NewV4()
	key := auther.StdProdAuther.GetAPIKey(id)

	fmt.Println(id)
	fmt.Println(key)
}
