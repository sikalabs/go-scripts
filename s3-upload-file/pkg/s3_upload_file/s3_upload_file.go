package s3_upload_file

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Options struct {
	FilePath  string
	Bucket    string
	Key       string
	Region    string
	Endpoint  string
	PathStyle bool
	ACL       string
	AccessKey string
	SecretKey string
}

func UploadFile(opts Options) error {
	file, err := os.Open(opts.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", opts.FilePath, err)
	}
	defer file.Close()

	ctx := context.Background()

	configOpts := []func(*config.LoadOptions) error{
		config.WithRegion(opts.Region),
	}
	if opts.AccessKey != "" && opts.SecretKey != "" {
		configOpts = append(configOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(opts.AccessKey, opts.SecretKey, ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if opts.Endpoint != "" {
			o.BaseEndpoint = aws.String(opts.Endpoint)
		}
		o.UsePathStyle = opts.PathStyle
	})

	uploader := manager.NewUploader(client)

	input := &s3.PutObjectInput{
		Bucket: aws.String(opts.Bucket),
		Key:    aws.String(opts.Key),
		Body:   file,
	}
	if opts.ACL != "" {
		input.ACL = types.ObjectCannedACL(opts.ACL)
	}

	_, err = uploader.Upload(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file to s3://%s/%s: %w", opts.Bucket, opts.Key, err)
	}

	return nil
}
