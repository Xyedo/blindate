package s3

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/xyedo/blindate/pkg/attachment/dtos"
)

func NewS3(bucketName string) *attachmentManager {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Panicln(err)
	}
	client := s3.NewFromConfig(cfg)
	return &attachmentManager{
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = 10 << 20
		}),
		presignClient: s3.NewPresignClient(client),
		s3client:      client,
		bucketName:    bucketName,
	}
}

type attachmentManager struct {
	uploader      *manager.Uploader
	presignClient *s3.PresignClient
	s3client      *s3.Client
	bucketName    string
}

func (a *attachmentManager) UploadBlob(file io.Reader, attachment dtos.Uploader) (string, error) {
	//TODO: better error handling
	key := attachment.Prefix + "/" + uuid.New().String() + attachment.Ext

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(a.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: attachment.Length,
		ContentType:   aws.String(attachment.ContentType),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (a *attachmentManager) DeleteBlob(key string) error {
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

func (a *attachmentManager) GetPresignedUrl(key string) (string, error) {
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
