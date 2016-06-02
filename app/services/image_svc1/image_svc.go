package image_svc1

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/bulk_transporter_svc"
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
