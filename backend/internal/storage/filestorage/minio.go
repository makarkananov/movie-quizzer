package filestorage

import (
	"context"
	"fmt"

	"movie-quizzer/backend/internal/service"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	client *minio.Client
	bucket string
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinIO, error) {
	c, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIO{
		client: c,
		bucket: bucket,
	}, nil
}

func (m *MinIO) GetMedia(bucket, file string) (service.MediaStream, error) {
	// Используем переданный bucket, если он указан, иначе используем дефолтный
	targetBucket := bucket
	if targetBucket == "" {
		targetBucket = m.bucket
	}

	// Сначала проверяем, существует ли объект
	ctx := context.Background()
	_, err := m.client.StatObject(ctx, targetBucket, file, minio.StatObjectOptions{})
	if err != nil {
		return service.MediaStream{}, fmt.Errorf("object %s not found in bucket %s: %w", file, targetBucket, err)
	}

	obj, err := m.client.GetObject(
		ctx,
		targetBucket,
		file,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return service.MediaStream{}, fmt.Errorf("failed to get object %s from bucket %s: %w", file, targetBucket, err)
	}

	stat, err := obj.Stat()
	if err != nil {
		obj.Close()
		return service.MediaStream{}, fmt.Errorf("failed to stat object %s: %w", file, err)
	}

	if stat.Size == 0 {
		obj.Close()
		return service.MediaStream{}, fmt.Errorf("file is empty or not found: %s", file)
	}

	// Определяем Content-Type если не указан
	contentType := stat.ContentType
	if contentType == "" {
		if len(file) > 4 {
			ext := file[len(file)-4:]
			switch ext {
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".png":
				contentType = "image/png"
			case ".mp4":
				contentType = "video/mp4"
			case ".webm":
				contentType = "video/webm"
			}
		}
	}

	return service.MediaStream{
		Reader:      obj,
		ContentType: contentType,
	}, nil
}
