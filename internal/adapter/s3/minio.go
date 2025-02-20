package s3

import (
	"context"
	"io"
	"strconv"

	"github.com/google/uuid"
	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/repository"
	"github.com/minio/minio-go/v7"
)

type MinIOS3RepositoryImpl struct {
	client *minio.Client
	cfg    *config.Config
}

func ProvideMinIOS3Repository(client *minio.Client, cfg *config.Config) repository.S3Repository {
	return &MinIOS3RepositoryImpl{
		client: client,
		cfg:    cfg,
	}
}

func (r *MinIOS3RepositoryImpl) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string, folder string) (minio.UploadInfo, error) {
	limitSize, err := strconv.ParseInt(r.cfg.FileUploadSizeLimitMB, 0, 64)
	if err != nil {
		return minio.UploadInfo{}, exception.ErrTypeConversion
	}
	if size > limitSize<<20 {
		return minio.UploadInfo{}, exception.ErrFileSizeLimit
	}

	_, err = r.client.StatObject(ctx, r.cfg.S3BucketName, objectName, minio.StatObjectOptions{})
	if err == nil {
		return minio.UploadInfo{}, exception.ErrFileAlreadyExists
	}

	objectName = folder + "/" + uuid.New().String() + "_" + objectName

	uploadInfo, err := r.client.PutObject(ctx, r.cfg.S3BucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	return uploadInfo, err
}
