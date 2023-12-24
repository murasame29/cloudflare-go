package client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type awsConnect struct {
	accountID       string
	endpoint        string
	accessKeyID     string
	accessKeySecret string
}

type Connection interface {
	Connect(ctx context.Context) (*s3.Client, error)
}

func New(
	accountID,
	endpoint,
	accessKeyID,
	accessKeySecret string,
) Connection {
	return &awsConnect{
		accountID:       accountID,
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
	}
}

func (a *awsConnect) newResolver() aws.EndpointResolverWithOptions {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               fmt.Sprintf("https://%s.%s", a.accountID, a.endpoint),
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom,
		}, nil
	})
}

func (a *awsConnect) Connect(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(a.newResolver()),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.accessKeyID,
				a.accessKeySecret,
				"",
			)),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}
