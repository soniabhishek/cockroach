package image_svc1

import (
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories/batch_process_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/bulk_transporter_svc"
)

//Transforms imager to a carrier. The imageCarrier knows how to transport an image.  A list of carriers can be passed to a bulkTransporter
type imageCarrier struct {
	imager
}

func (i *imageCarrier) Pick() error {
	return i.Download()
}

func (i *imageCarrier) Drop() error {
	return i.Upload()
}

func (i *imageCarrier) Run() error {
	return nil
}

func (i *imageCarrier) GetId() uuid.UUID {
	return i.ImageDictionaryNew.ID
}

var _ bulk_transporter_svc.ICarrier = &imageCarrier{}

//--------------------------------------------------------------------------------//

type imageCarrierLogger struct {
	batchProcessRepo batch_process_repo.IBatchProcessRepo
	batchProcess     models.BatchProcess
	totalImages      int
	successCounter   int
	failureCounter   int
}

func (l *imageCarrierLogger) LogStartTransport() {
	fmt.Println("image download started")
}
func (l *imageCarrierLogger) LogPickSuccess(c bulk_transporter_svc.ICarrier) {
	//ic := c.(*imageCarrier)
	//fmt.Println("PickSuccess", ic.OriginalUrl)
}
func (l *imageCarrierLogger) LogPickFailure(c bulk_transporter_svc.ICarrier) {
	//ic := c.(*imageCarrier)
	//fmt.Println("PickFailure", ic.OriginalUrl)
	l.failureCounter++
}
func (l *imageCarrierLogger) LogDropSuccess(c bulk_transporter_svc.ICarrier) {
	//ic := c.(*imageCarrier)
	//fmt.Println("DropSuccess", ic.OriginalUrl)
	l.successCounter++
	l.updateStats()
}
func (l *imageCarrierLogger) LogDropFailure(c bulk_transporter_svc.ICarrier) {
	//ic := c.(*imageCarrier)
	//fmt.Println("DropFailure", ic.OriginalUrl)
	l.failureCounter++
}
func (l *imageCarrierLogger) LogFinishTransport() {
	fmt.Println("image download finished")
}

func (l *imageCarrierLogger) updateStats() {
	l.batchProcess.Completion = float64(l.successCounter) * 100 / float64(l.totalImages)
	fmt.Println(l.batchProcess.Completion)
}
