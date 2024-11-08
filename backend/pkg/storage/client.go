package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Bucket interface {
	Upload(ctx context.Context, key string, body io.Reader, opts ...ObjectOption) error
	GetPublicURL(ctx context.Context, key string) string
	GetPresignedURL(ctx context.Context, key string, expireDuration time.Duration) (string, error)
}

type bucket struct {
	name    string
	region  string
	storage *s3.Client
}

type Params struct {
	BucketName      string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type ObjectOption func(*s3.PutObjectInput)

func WithACL(acl types.ObjectCannedACL) ObjectOption {
	return func(opts *s3.PutObjectInput) {
		opts.ACL = acl
	}
}

func WithCacheControl(maxAge time.Duration) ObjectOption {
	const format = "max-age=%.0f,public"
	return func(attr *s3.PutObjectInput) {
		attr.CacheControl = aws.String(fmt.Sprintf(format, maxAge.Seconds()))
	}
}

func NewClient(ctx context.Context, params *Params) (Bucket, error) {
	creds := credentials.NewStaticCredentialsProvider(params.AccessKeyID, params.SecretAccessKey, "")
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(params.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	storage := s3.NewFromConfig(cfg)
	return &bucket{
		name:    params.BucketName,
		region:  params.Region,
		storage: storage,
	}, nil
}

func (c *bucket) Upload(ctx context.Context, key string, body io.Reader, opts ...ObjectOption) error {
	objectInput := &s3.PutObjectInput{
		Bucket: aws.String(c.name),
		Key:    aws.String(key),
		Body:   body,
	}
	for i := range opts {
		opts[i](objectInput)
	}
	_, err := c.storage.PutObject(ctx, objectInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *bucket) GetPublicURL(ctx context.Context, key string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("s3.%s.amazonaws.com", c.region),
		Path:   fmt.Sprintf("/%s/%s", c.name, key),
	}
	return u.String()
}

func (c *bucket) GetPresignedURL(ctx context.Context, key string, expireDuration time.Duration) (string, error) {
	presign := s3.NewPresignClient(c.storage)

	request, err := presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.name),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireDuration))
	if err != nil {
		return "", nil
	}

	return request.URL, nil
}
