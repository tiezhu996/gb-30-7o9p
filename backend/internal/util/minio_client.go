package util

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gbadopt/gbadopt/internal/config"
)

// MinIOClient wraps the minio-go client for object storage operations.
type MinIOClient struct {
	client *minio.Client
	bucket string
}

// NewMinIOClient connects to MinIO and ensures the bucket exists.
func NewMinIOClient(cfg *config.Config) (*MinIOClient, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOUser, cfg.MinIOPassword, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket check: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio make bucket: %w", err)
		}
	}
	return &MinIOClient{client: client, bucket: cfg.MinIOBucket}, nil
}

// Upload stores a reader under an object name and returns the object key.
func (m *MinIOClient) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio put object: %w", err)
	}
	return objectName, nil
}

// Get streams an object's contents.
func (m *MinIOClient) Get(ctx context.Context, objectName string) (io.ReadCloser, int64, string, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, "", fmt.Errorf("minio get object: %w", err)
	}
	stat, err := obj.Stat()
	if err != nil {
		return nil, 0, "", fmt.Errorf("minio stat object: %w", err)
	}
	return obj, stat.Size, stat.ContentType, nil
}

// Remove deletes an object.
func (m *MinIOClient) Remove(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}
