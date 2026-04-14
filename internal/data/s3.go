package data

import (
	"context"
	"file/internal/biz"
	"file/internal/conf"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-kratos/kratos/v2/log"
)

type s3StorageProvider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	log           *log.Helper
}

func NewS3StorageProvider(c *conf.Data, logger log.Logger) (biz.StorageProvider, error) {
	if c.S3 == nil {
		return nil, fmt.Errorf("s3 configuration is missing")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.S3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.S3.AccessKey, c.S3.SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if c.S3.Endpoint != "" {
			o.BaseEndpoint = aws.String(c.S3.Endpoint)
		}
	})
	presignClient := s3.NewPresignClient(client)

	return &s3StorageProvider{
		client:        client,
		presignClient: presignClient,
		bucket:        c.S3.Bucket,
		log:           log.NewHelper(logger),
	}, nil
}

func (s *s3StorageProvider) GetUploadUrl(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error) {
	req, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *s3StorageProvider) GetDownloadUrl(ctx context.Context, key string, filename string, expiry time.Duration) (string, error) {
	disposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(s.bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(disposition),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *s3StorageProvider) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
