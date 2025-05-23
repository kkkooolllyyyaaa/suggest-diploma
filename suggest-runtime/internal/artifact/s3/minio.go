package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

type Minio struct {
	client     *minio.Client
	bucketName string
}

func NewMinio(cfg Config) (*Minio, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Minio: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil || !exists {
		return nil, fmt.Errorf("couldn't find bucket %s: %w", cfg.BucketName, err)
	}

	return &Minio{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

// DownloadFileToFile скачивает файл из MinIO/S3 и сохраняет в локальный файл
func (m *Minio) DownloadFileToFile(objectName, filePath string) error {
	ctx := context.Background()
	err := m.client.FGetObject(ctx, m.bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download file from %s to %s: %w", objectName, filePath, err)
	}
	return nil
}
