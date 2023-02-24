package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewS3Client creates a new minio.Client
func NewS3Client(endpoint string, accessKeyID string, secretAccessKey string, secure bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
