package client

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	once     sync.Once
)

// InitializeS3Client initializes the S3 client and sets it to a package-level variable.
func InitializeS3Client() {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		s3Client = s3.NewFromConfig(cfg)
	})
}

// GetS3Client returns the initialized S3 client.
func GetS3Client() *s3.Client {
	if s3Client == nil {
		InitializeS3Client()
	}
	return s3Client
}
