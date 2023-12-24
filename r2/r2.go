package r2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2crud interface {
	GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error)
	ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error)

	PublishPresignedObjectURL(ctx context.Context, key string) (string, error)
	ListObjectsURL(ctx context.Context, prefix string) ([]string, error)

	UploadObject(ctx context.Context, file []byte, key string) error
	DeleteObject(ctx context.Context, key string) error
}

type r2crud struct {
	client        *s3.Client
	PresignClient *s3.PresignClient
	bucket        string
	Config        Config
}

type Config struct {
	PresignLinkExpired time.Duration
}

func NewR2CRUD(bucket string, client *s3.Client, presignLinkExpired int) R2crud {
	return &r2crud{
		bucket:        bucket,
		client:        client,
		PresignClient: s3.NewPresignClient(client),
		Config: Config{
			PresignLinkExpired: time.Duration(presignLinkExpired),
		},
	}
}
