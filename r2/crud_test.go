package r2_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	env "github.com/murasame29/cloudflare-go/config"
	"github.com/murasame29/cloudflare-go/r2"
	"github.com/stretchr/testify/require"
)

type awsConfig struct {
	accountID        string
	accountKeyID     string
	accountKeySecret string
	bucketName       string
	s3Endpoint       string
}

func NewAwsConfig() *awsConfig {
	return &awsConfig{
		accountID:        env.Config.CfAccountID,
		accountKeyID:     env.Config.CfAccountKeyID,
		accountKeySecret: env.Config.CfAccountKeySecret,
		bucketName:       env.Config.CfBucketName,
		s3Endpoint:       env.Config.CfS3Endpoint,
	}
}

func (a *awsConfig) Client() (*s3.Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", a.accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(a.accountKeyID, a.accountKeySecret, "")),
	)
	return s3.NewFromConfig(cfg), err
}
func TestR2(t *testing.T) {
	image, err := os.Open("../images.png")
	require.NoError(t, err)
	defer image.Close()

	client, err := NewAwsConfig().Client()
	require.NoError(t, err)

	r2 := r2.NewR2crud(client)

	output, err := r2.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(env.Config.CfBucketName),
		Key:         aws.String("hoge.png"),
		Body:        image,
		ContentType: aws.String("image/png"),
	})

	require.NoError(t, err)
	require.NotNil(t, output)

}
