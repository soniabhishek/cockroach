package image_svc

import (
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.com/playment-main/angel/app/config"
)

func uploadToS3(img io.ReadCloser) {

	defer img.Close()

	awsConfig := aws.NewConfig().
		WithRegion(config.Get(config.AWS_REGION)).
		WithCredentials(credentials.NewStaticCredentials(config.Get(config.AWS_ACCESS_ID), config.Get(config.AWS_SECRET_KEY), ""))

	uploader := s3manager.NewUploader(session.New(awsConfig))

	uploader.Concurrency = 60

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   img,
		Bucket: aws.String("playmentdevelopment"),
		Key:    aws.String("default/test1/111.png"),
	})

	if err != nil {
		log.Fatalln("Failed to upload", err)
	}
	log.Println("Uploaded", result.Location)
}
