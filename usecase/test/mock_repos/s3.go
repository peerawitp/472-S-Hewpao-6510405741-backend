package mock_repos

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
)

type MockS3Repository struct {
	mock.Mock
}

func (m *MockS3Repository) UploadFile(ctx context.Context, objectName string, content io.Reader, size int64, contentType string, folder string) (minio.UploadInfo, error) {
	args := m.Called(ctx, objectName, content, size, contentType, folder)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockS3Repository) GetSignedURL(ctx context.Context, bucketName string, objectName string, expires time.Duration) (string, error) {
	args := m.Called(ctx, bucketName, objectName, expires)
	return args.String(0), args.Error(1)
}
