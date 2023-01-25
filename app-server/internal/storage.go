package internal

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage interface {
	CurrentBucket() string
	PresignedUrl(method, key string, hours int64) (string, error)
	FileExist(key string) bool
	FileUpload(filepath, key string) error
}

type S3Storage struct {
	Bucket string
	client *s3.S3
}

func (storage *S3Storage) CurrentBucket() string {
	return storage.Bucket
}

func (storage *S3Storage) PresignedUrl(method, key string, hours int64) (string, error) {
	var req *request.Request
	if method == "PUT" {
		req, _ = storage.client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(storage.Bucket),
			Key:    aws.String(key),
		})
	} else {
		req, _ = storage.client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(storage.Bucket),
			Key:    aws.String(key),
		})
	}

	return req.Presign(time.Duration(hours) * time.Hour)
}

func (storage *S3Storage) FileExist(key string) bool {
	_, err := storage.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(storage.Bucket),
		Key:    aws.String(key),
	})
	return err == nil
}

func (storage *S3Storage) FileUpload(filepath, key string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	info, _ := f.Stat()
	var fileSize int64 = info.Size()
	fileBuffer := make([]byte, fileSize)
	f.Read(fileBuffer)

	_, err = storage.client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(storage.Bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func NewStorage(configs *Configs) *S3Storage {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(configs.Region)},
	)
	if err != nil {
		panic(err)
	}
	client := s3.New(sess)

	return &S3Storage{client: client, Bucket: configs.StorageBucket}
}
