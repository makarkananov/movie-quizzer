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
	obj, err := m.client.GetObject(
		context.Background(),
		m.bucket,
		file,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return service.MediaStream{}, err
	}

	stat, err := obj.Stat()
	if err != nil {
		return service.MediaStream{}, err
	}

	if stat.Size == 0 {
		return service.MediaStream{}, fmt.Errorf("file is empty or not found")
	}

	return service.MediaStream{
		Reader:      obj,
		ContentType: stat.ContentType,
	}, nil
}
