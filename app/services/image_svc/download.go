package image_svc

import (
	"io/ioutil"
	"net/http"

	"github.com/crowdflux/angel/app/models"
)

type image struct {
	models.ImageContainer
	downloadedImage
}

type downloadedImage struct {
	bty []byte
	err error
}

func DownloadImagesFromList(ics []models.ImageContainer) error {

	var channelSize = 60

	downLoadCh := make(chan downloadedImage, channelSize)
	//uploadCh := make(chan int, channelSize)

	//start := time.Now()

	for i := 0; i < len(ics)+channelSize; i++ {

		if i < len(ics) {
			go download2(ics[i].Url, downLoadCh)
		}
		if i >= channelSize {
			<-downLoadCh
			//imgBty := <-downLoadCh
			//go upload(imageContainer{imgBty, ics[i]}, uploadCh)
			//fmt.Println(i-channelSize, time.Since(start))
		}
	}
	return nil
}

func download2(url string, c chan downloadedImage) {

	resp, err := http.Get(url)
	//check(err)
	if err != nil {
		c <- downloadedImage{nil, err}
		return
	}
	defer resp.Body.Close()

	bty, err := ioutil.ReadAll(resp.Body)
	c <- downloadedImage{bty, err}
}

func upload(image image, c chan int) {
	//
	//bty, err := ioutil.ReadAll(image)
	////check(err)
	//if err != nil {
	//	log.Fatal("Error in reading downloaded image: ", err)
	//	return
	//}
	//
	//f, err := os.Create("temp/" + image.Id + ".jpeg")
	////check(err)
	//if err != nil {
	//	log.Fatal("Error in creating image: ", err)
	//	return
	//}
	//
	//defer f.Close()
	//
	//_, err = f.Write(bty)
	//
	//if err != nil {
	//	log.Fatal("Error in saving image: ", err)
	//	return
	//}
}
