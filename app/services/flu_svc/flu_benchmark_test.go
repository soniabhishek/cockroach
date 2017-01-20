package flu_svc

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkFluUpload(b *testing.B) {
	path := "./to_soni.csv"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Println("running")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload", filepath.Base(path))
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}
	client := http.DefaultClient
	req, err := http.NewRequest("POST", "http://localhost:8999/api/v0/project/9e9f84e2-1ab7-4118-bd73-51af7a9ded61/csv/feedline", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := client.Do(req)
	response.Body.Close()
	fmt.Println(response, err)

}
