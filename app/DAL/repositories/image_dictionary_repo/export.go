package image_dictionary_repo

import (
	"io"
	"net/http"

	"gitlab.com/playment-main/angel/app/models"
)

type IImageDictionaryRepo interface {
	UploadToS3(models.ImageDictionary) error
	Download(link string) (i io.ReadCloser, err error)
}

func New() IImageDictionaryRepo {
	return &imageDictionaryRepo{}
}

type imageDictionaryRepo struct {
}

func (r *imageDictionaryRepo) UploadToS3(m models.ImageDictionary) error {
	return nil
}

func (r *imageDictionaryRepo) Download(link string) (i io.ReadCloser, err error) {
	resp, err := http.Get(link)
	//check(err)
	if err != nil {
		return
	}
	i = resp.Body
	return
}
