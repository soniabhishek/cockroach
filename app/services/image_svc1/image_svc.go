package image_svc1

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/bulk_transporter_svc"
)

type imageService struct {
}

func (s *imageService) BulkDownloadImages(imgs []models.ImageDictionaryNew) (batchId uuid.UUID, err error) {

	var imgCarriers []bulk_transporter_svc.ICarrier = make([]bulk_transporter_svc.ICarrier, len(imgs))
	for i, img := range imgs {
		imgCarriers[i] = &imageCarrier{imager{ImageDictionaryNew: img}}
	}

	imageLogger := &imageCarrierLogger{
		totalImages: len(imgCarriers),
	}

	b := bulk_transporter_svc.New(imageLogger)

	err = b.BulkTransport(imgCarriers)

	return uuid.Nil, nil
}
