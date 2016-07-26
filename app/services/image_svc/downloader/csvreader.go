package downloader

import (
	"encoding/csv"
	"os"

	"github.com/crowdflux/angel/app/models"
)

func ReadFromTempDir(fileName string) (images []models.ImageContainer, err error) {

	file, err := os.Open("temp/" + fileName)

	if err != nil {
		return
	}

	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		return
	}

	records = records[1:len(records)]

	for _, record := range records {
		images = append(images, models.ImageContainer{
			Id:  record[0],
			Url: record[1],
		})
	}
	return images, nil

}
