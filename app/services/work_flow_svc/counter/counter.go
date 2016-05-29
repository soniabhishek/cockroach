package counter

import (
	"fmt"
	"gitlab.com/playment-main/angel/app/models"
)

func Print(flu models.FeedLineUnit) {
	fmt.Println(flu.ID.String() + " reached in " + flu.Step)
}
