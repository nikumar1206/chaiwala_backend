package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	client     *s3.Client
	BucketName string
}

// New initializes the client with credentials from env
func New(ctx context.Context, awsRegion, s3BucketName string) S3Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
	)
	if err != nil {
		panic(err)
	}

	return S3Client{
		client:     s3.NewFromConfig(cfg),
		BucketName: s3BucketName,
	}
}

// Upload uploads an image to the bucket
func (s *S3Client) Upload(ctx context.Context, key string, fileReader io.Reader, contentType string) error {
	uploader := manager.NewUploader(s.client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(key),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})

	return err
}

// Upload uploads an image to the bucket
func (s *S3Client) UploadV2(ctx context.Context, key string, content []byte, contentType string) error {
	_, err := s.client.PutObject(
		ctx, &s3.PutObjectInput{
			Bucket:               aws.String(s.BucketName),
			Key:                  aws.String(key),
			Body:                 bytes.NewReader(content),
			ContentType:          aws.String(contentType),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		})

	return err
}

// Download fetches an object and returns its bytes
func (s *S3Client) Download(ctx context.Context, key string) ([]byte, error) {
	downloader := manager.NewDownloader(s.client)

	buf := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Update is just a re-upload
func (s *S3Client) Update(ctx context.Context, key string, fileReader io.Reader, contentType string) error {
	return s.Upload(ctx, key, fileReader, contentType)
}

// Delete removes the object from the bucket
func (s *S3Client) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return err
}
