package services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/byt3roof/odoo-backup/internal/conf"
)

func UploadToS3(ctx context.Context, filePath string, key string) error {
	env, err := conf.LoadConfig()
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(
			env.AccessKey,
			env.SecretKey,
			"",
		),
	))
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(env.S3Endpoint)
		o.Region = "auto"
	})

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &env.Bucket,
		Key:    &key,
		Body:   file,
	})

	return err
}
