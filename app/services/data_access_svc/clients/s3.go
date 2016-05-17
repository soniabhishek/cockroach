package clients

import (
	"io"
	"sync"

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
var onceS3 sync.Once

func init() {
	onceS3.Do(func() {
		uploader = initS3()
	})
}

func initS3() *s3manager.Uploader {

	config.SetEnvironment(config.Development)

	awsConfig := aws.NewConfig().
		WithRegion(config.GetVal(config.AWS_REGION)).
		WithCredentials(credentials.NewStaticCredentials(config.GetVal(config.AWS_ACCESS_ID), config.GetVal(config.AWS_SECRET_KEY), ""))

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
