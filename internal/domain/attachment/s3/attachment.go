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
	"github.com/xyedo/blindate/internal/common/ids"
	"github.com/xyedo/blindate/internal/domain/attachment/entities"
	"github.com/xyedo/blindate/internal/infrastructure"
)

var Manager *attachmentManager

func InitS3Manager() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Panicln(err)
	}
	client := s3.NewFromConfig(cfg)
	Manager = &attachmentManager{
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = 10 << 20
		}),
		presignClient: s3.NewPresignClient(client),
		s3client:      client,
		bucketName:    infrastructure.Config.AWS.BucketName,
	}

}

type attachmentManager struct {
	uploader      *manager.Uploader
	presignClient *s3.PresignClient
	s3client      *s3.Client
	bucketName    string
}

func (a *attachmentManager) UploadAttachment(ctx context.Context, file io.Reader, attachment entities.Attachment) (string, error) {
	//TODO: better error handling
	key := attachment.Prefix + "/" + ids.Attachment() + attachment.Ext

	_, err := a.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(attachment.ContentType),
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (a *attachmentManager) DeleteBlob(ctx context.Context, key string) error {
	_, err := a.s3client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Key:    &key,
		Bucket: aws.String(a.bucketName),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *attachmentManager) GetPresignedUrl(ctx context.Context, key string, expires time.Duration) (string, error) {
	presignRes, err := a.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return presignRes.URL, nil
}
