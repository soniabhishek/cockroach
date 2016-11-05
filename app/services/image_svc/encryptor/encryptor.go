package encryptor

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc"
)

func GetEncryptedUrls(imageField interface{}) (urlSlice []string, err error) {

	encResult, err := clients.GetLuigiClient().GetEncryptedUrls(imageField)
	var d = 0
	for item := range encResult {
		if item["valid"] == false {
			err = flu_svc.ErrDataMissing
			plog.Error("Image Encryption step : Image not encryptable", err)
			return
		}
		urlSlice[d] = string(item["playment_url"])
		d++
	}
	return
}
