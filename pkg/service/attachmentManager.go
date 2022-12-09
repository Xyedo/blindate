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
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
)

type Attachment interface {
	UploadBlob(file io.Reader, attach domain.Uploader) (string, error)
	DeleteBlob(key string) error
	GetPresignedUrl(key string) (string, error)
}

func NewS3(bucketName string) *attachment {
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
		s3client:      client,
		bucketName:    bucketName,
	}
}

type attachment struct {
	uploader      *manager.Uploader
	presignClient *s3.PresignClient
	s3client      *s3.Client
	bucketName    string
}

func (a *attachment) UploadBlob(file io.Reader, attach domain.Uploader) (string, error) {
	//TODO: better error handling
	key := attach.Prefix + "/" + util.RandomUUID() + attach.Ext

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(a.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: attach.Length,
		ContentType:   aws.String(attach.ContentType),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (a *attachment) DeleteBlob(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Key:    &key,
		Bucket: aws.String(a.bucketName),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *attachment) GetPresignedUrl(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	presignRes, err := a.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	})
	if err != nil {
		return "", err
	}
	return presignRes.URL, nil
}
