package r2

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2crud interface {
	ReadAll(context.Context, string) ([]*ReadAllOutput, error)
	Upload(context.Context, *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	Download(context.Context, *s3.GetObjectInput) ([]byte, error)
	Delete(context.Context, *s3.DeleteObjectInput) error
	SetOption(func(*s3.Options))
}

type r2Crud struct {
	c      *s3.Client
	option []func(*s3.Options)
}

func NewR2crud(c *s3.Client) R2crud {
	return &r2Crud{
		c: c,
	}
}

type ReadAllOutput struct {
	ChecksumAlgorithm string `json:"ChecksumAlgorithm"`
	ETag              string `json:"ETag"`
	Key               string `json:"Key"`
	LastModified      string `json:"LastModified"`
	Owner             string `json:"Owner"`
	Size              int    `json:"Size"`
	StorageClass      string `json:"StorageClass"`
}

func (r *r2Crud) ReadAll(ctx context.Context, bucketName string) ([]*ReadAllOutput, error) {
	listObjectsOutput, err := r.c.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, err
	}

	var results []*ReadAllOutput
	for _, object := range listObjectsOutput.Contents {
		var result *ReadAllOutput
		obj, err := json.MarshalIndent(object, "", "\t")
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(obj, &result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *r2Crud) Upload(ctx context.Context, arg *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	output, err := r.c.PutObject(ctx, arg)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (r *r2Crud) Download(ctx context.Context, arg *s3.GetObjectInput) ([]byte, error) {
	image, err := r.c.GetObject(ctx, arg)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(image.Body)
	defer image.Body.Close()

	return buf.Bytes(), nil
}

func (r *r2Crud) Delete(ctx context.Context, arg *s3.DeleteObjectInput) error {
	_, err := r.c.DeleteObject(ctx, arg)
	return err
}

func (r *r2Crud) SetOption(f func(*s3.Options)) {
	r.option = append(r.option, f)
}
