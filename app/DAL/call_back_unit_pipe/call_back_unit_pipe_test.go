package call_back_unit_pipe

import (
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	//	"github.com/stretchr/testify/assert"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/feed_line"
)

func TestNew(t *testing.T) {

	CbuQ := New("test1")

	fluid1 := uuid.FromStringOrNil("18846509-8ec6-49f6-8312-7d5d176fc0ac")
	refid1 := "38a27427-9474-4e76-bd77-bc7dee2463e5"

	fluid2 := uuid.FromStringOrNil("a2f2e617-25b5-412c-af06-737c801ebb99")
	refid2 := "44a27427-9474-4e76-bd77-bc7dee2463e5"
	var fluout = []models.FluOutputStruct{{
		ID:          fluid1,
		ReferenceId: refid1,
		Tag:         "sometag",
		Status:      "somestatus",
		Result:      "someresult",
	},
		{
			ID:          fluid2,
			ReferenceId: refid2,
			Tag:         "sometag2",
			Status:      "somestatus2",
			Result:      "someresult2",
		},
	}

	var flusSent = make(map[uuid.UUID]feed_line.FLU)
	flusSent[fluid1] = feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:          fluid1,
			ReferenceId: refid1,
		},
	}
	flusSent[fluid2] = feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:          fluid2,
			ReferenceId: refid2,
		},
	}
	pConfig := models.ProjectConfiguration{
		ProjectId:   uuid.FromStringOrNil("38a27427-9474-4e76-bd77-bc7dee2463e5"),
		PostBackUrl: "someurl.insomedomain.com/post",
	}
	CbuQ.Push(CBU{
		FluOutputObj:  fluout,
		FlusSent:      flusSent,
		ProjectConfig: pConfig,
		RetryLeft:     3,
	})

	cbu := <-CbuQ.Receiver()

	fmt.Println(cbu.ProjectConfig)
	fmt.Println(cbu.FluOutputObj)
	fmt.Print(cbu.FlusSent)
	cbu.ConfirmReceive()
	//assert.EqualValues(t, fluId, flu.)

}
