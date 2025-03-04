package repository

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type S3Repository interface {
	UploadFile(ctx context.Context, objectName string, content io.Reader, size int64, contentType string, folder string) (minio.UploadInfo, error)
	GetSignedURL(ctx context.Context, bucketName string, objectName string, expires time.Duration) (string, error)
}
