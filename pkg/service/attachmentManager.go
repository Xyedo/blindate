package service

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/xyedo/blindate/pkg/util"
)

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
		downloader: manager.NewDownloader(client),
	}
}

type attachment struct {
	uploader   *manager.Uploader
	downloader *manager.Downloader
}

func (a *attachment) UploadBlob(file []byte, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := util.RandomUUID()
	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("blindate-bucket"),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}
