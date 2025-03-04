package util

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/domain/exception"
	"github.com/hewpao/hewpao-backend/repository"
)

func FileManage(fileHeaders *multipart.Form, field string, limit int) ([]io.Reader, []*multipart.FileHeader, error) {
	files := fileHeaders.File[field]
	if len(files) == 0 {
		return nil, nil, exception.ErrFileIsNull
	}
	if len(files) > limit {
		return nil, nil, exception.ErrFileCountLimit
	}

	var fileReaders []io.Reader

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, nil, err
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, nil, err
		}
		reader := bytes.NewReader(content)
		fileReaders = append(fileReaders, reader)

		defer file.Close()
	}

	return fileReaders, files, nil
}

func GetUrls(minioRepo repository.S3Repository, ctx context.Context, cfg *config.Config, images []string) ([]string, error) {
	duration, err := time.ParseDuration(cfg.S3Expiration)
	urls := []string{}
	if err != nil {
		return nil, err
	}

	for _, img := range images {
		path := strings.SplitN(img, "hewpao-s3/", 2)
		url, err := minioRepo.GetSignedURL(ctx, cfg.S3BucketName, path[1], duration)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
