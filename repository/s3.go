package repository

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type S3Repository interface {
	UploadFile(ctx context.Context, objectName string, content io.Reader, size int64, contentType string, folder string) (minio.UploadInfo, error)
}
