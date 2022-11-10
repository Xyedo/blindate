package service

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/xyedo/blindate/pkg/util"
)

const BUCKET_NAME = "blindate-bucket"

func NewS3() *attachment {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Panicln(err)
	}
	client := s3.NewFromConfig(cfg)
	return &attachment{
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = 10 << 20
		}),
		presignClient: s3.NewPresignClient(client),
	}
}

type attachment struct {
	uploader      *manager.Uploader
	presignClient *s3.PresignClient
}

func (a *attachment) UploadBlob(file io.Reader, length int64, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := util.RandomUUID()
	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(BUCKET_NAME),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: length,
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (a *attachment) GetPresignedUrl(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	presignRes, err := a.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	})
	if err != nil {
		return "", err
	}
	return presignRes.URL, nil
}
