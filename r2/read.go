package r2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (r *r2crud) ListObjectsURL(ctx context.Context, prefix string) ([]string, error) {
	object, err := r.ListObjects(ctx, prefix)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range object.Contents {
		accessURL, err := r.PublishPresignedObjectURL(ctx, *obj.Key)
		if err != nil {
			return nil, err
		}
		files = append(files, accessURL)
	}

	return files, nil
}

func (r *r2crud) GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	return r.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
}

func (r *r2crud) ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error) {
	return r.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(prefix),
	})
}

func (r *r2crud) PublishPresignedObjectURL(ctx context.Context, key string) (string, error) {
	object, err := r.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:          aws.String(r.bucket),
		Key:             aws.String(key),
		ResponseExpires: aws.Time(time.Now().Add(r.Config.PresignLinkExpired * time.Hour)),
	})
	if err != nil {
		return "", err
	}

	return object.URL, nil
}
