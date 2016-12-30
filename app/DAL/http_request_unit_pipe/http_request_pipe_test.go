package http_request_pipe

/*

import (
	"testing"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNew(t *testing.T) {

	cr := New("test1")

	fmcrId := uuid.NewV4()

	cr.Push(FMCR{
		CallBack: models.Flu_output_struct{
		},
	})

	fmcr := <-cr.Receiver()

	fmcr.ConfirmReceive()

	assert.EqualValues(t, fmcr, fmcr.ID)

}

func TestFeedline_Load(t *testing.T) {
	cr := New("test12")

	fmcrs := cr.Receiver()

	go func() {

		for {
			cr.Push(FMCR{
				Flu_output_struct: models.Flu_output_struct{
					ID: uuid.NewV4(),
				},
			})

		}
	}()

	go func() {

		for {
			fmcr := <-fmcrs
			fmcr.ConfirmReceive()

		}
	}()

	time.Sleep(time.Duration(1) * time.Second)

}
*/
