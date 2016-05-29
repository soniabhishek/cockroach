package clients

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.com/playment-main/angel/app/config"
)

//type IS3Client interface {
//	Upload(data io.Reader, s3Bucket string, key string) error
//	Download(link string) (i io.ReadCloser, err error)
//}

var uploader *s3manager.Uploader

func init() {
	uploader = initS3()
}

func initS3() *s3manager.Uploader {

	awsConfig := aws.NewConfig().
		WithRegion(config.Get(config.AWS_REGION)).
		WithCredentials(credentials.NewStaticCredentials(config.Get(config.AWS_ACCESS_ID), config.Get(config.AWS_SECRET_KEY), ""))

	uploader = s3manager.NewUploader(session.New(awsConfig))
	uploader.Concurrency = 60
	return uploader
}

func GetS3Client() s3Client {
	return s3Client{}
}

type s3Client struct {
}

func (s3Client) Upload(data io.Reader, s3Bucket string, key string) error {

	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:        data,
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String(key),
		ContentType: aws.String("image/jpeg"),
	})

	return err
}

func (s3Client) Download(link string) error {
	panic("not implemented")
	return nil
}
