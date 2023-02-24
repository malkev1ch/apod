package repository

import (
	"context"
	"fmt"

	"github.com/malkev1ch/apod/internal/model"
	"github.com/minio/minio-go/v7"
)

type FileRepository struct {
	client   *minio.Client
	endpoint string
}

// NewFileRepository creates new FileRepository.
func NewFileRepository(client *minio.Client, endpoint string) *FileRepository {
	return &FileRepository{client: client, endpoint: endpoint}
}

func (f *FileRepository) PutObject(ctx context.Context, bucket string, input *model.File) (string, error) {
	options := minio.PutObjectOptions{
		ContentType: input.ContentType,
		// grant access policy permissions
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	uploadInfo, err := f.client.PutObject(
		ctx,
		bucket,
		input.Name,
		input.File,
		input.Size,
		options,
	)
	if err != nil {
		return "", fmt.Errorf("f.client.PutObject: %w", err)
	}

	return f.generateURL(bucket, uploadInfo.Key), nil
}

func (f *FileRepository) generateURL(bucket string, key string) string {
	return fmt.Sprintf("%s/%s/%s", f.endpoint, bucket, key)
}
