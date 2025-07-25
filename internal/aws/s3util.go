package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Util struct {
	Client *s3.Client
	Bucket string
}

func NewS3Util(ctx context.Context, bucket string, region string) (*S3Util, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}
	client := s3.NewFromConfig(cfg)
	return &S3Util{
		Client: client,
		Bucket: bucket,
	}, nil
}

// Generate a presigned GET URL
func (s *S3Util) PresignGetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	psClient := s3.NewPresignClient(s.Client)
	req, err := psClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// Generate a presigned PUT URL
func (s *S3Util) PresignPutURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	psClient := s3.NewPresignClient(s.Client)
	req, err := psClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPrivate,
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// Optional: Upload a file directly (not presigned)
func (s *S3Util) UploadFile(ctx context.Context, key string, body []byte) error {
	uploader := manager.NewUploader(s.Client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   manager.ReadSeekCloser(bytes.NewReader(body)),
	})
	return err
}

func (s *S3Util) StreamObject(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
