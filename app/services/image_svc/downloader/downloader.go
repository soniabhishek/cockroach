package downloader

import (
	"bufio"
	"bytes"
	"fmt"
	i "image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"time"

	"github.com/nfnt/resize"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/data_access_svc/clients"
	"io/ioutil"
)

//Pass the imageContainerList and it will download those
//images in 60 concurrent threads
func DownloadFromList(imageList []models.ImageContainer, folderName string) error {

	if _, err := os.Stat("./temp/" + folderName); os.IsNotExist(err) {
		if err = os.Mkdir("./temp/"+folderName, os.ModePerm); err != nil {
			return err
		}
		if err = os.Mkdir("./temp/"+folderName+"/h200", os.ModePerm); err != nil {
			return err
		}
	}

	var channelSize = 60

	c := make(chan int, channelSize)

	start := time.Now()

	for i := 0; i < len(imageList)+channelSize; i++ {

		if i < len(imageList) {
			go download(imageList[i], folderName, c)
		}
		if i >= channelSize {
			<-c
			fmt.Println(i-channelSize, time.Since(start))
		}
	}
	return nil
}

// Find out if this implementation is correct.
// Is there any need to read the document (ioutil.ReadAll()) before
// saving it to the disk?
func download(image models.ImageContainer, folderName string, c chan int) {

	defer func() {
		c <- 1
	}()

	resp, err := http.Get(image.Url)
	//check(err)
	if err != nil {
		fmt.Println("Error in getting image: ", err)
		return
	}

	// Read the content
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	defer resp.Body.Close()

	err = clients.GetS3Client().Upload(bytes.NewReader(bodyBytes), "playmentproduction", "/public/"+folderName+"/"+image.Id+".jpeg")
	if err != nil {
		fmt.Println(err)
	}

	/*Mark it true to resize it to height 200px
	TODO This wont to go production until we figure out a better way
	to resize images as this takes lots of cpu & memory:
	On 60 concurrent downloads:
	Without this 35MB 4% CPU
	With this 1GB 30% CPU, time increases by 33%
	All of this is happening because of this line img, err := jpeg.Decode(b)
	*/
	if true {

		bt := bytes.NewReader(bodyBytes)

		img, _, err := i.Decode(bt) //This conversion takes lots of cpu & mem

		if err != nil {
			fmt.Println("error decoding the image", err)
			return
		}

		newImage := resize.Thumbnail(500, 500, img, resize.Lanczos3)

		var b bytes.Buffer
		byteWriter := bufio.NewReadWriter(bufio.NewReader(&b), bufio.NewWriter(&b))

		err = jpeg.Encode(byteWriter, newImage, nil)
		if err != nil {
			fmt.Println(err)
		}

		err = clients.GetS3Client().Upload(byteWriter, "playmentproduction", "/public/"+folderName+"/h200/"+image.Id+".jpeg")
		if err != nil {
			fmt.Println(err)
		}

	}

}
