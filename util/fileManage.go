package util

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/hewpao/hewpao-backend/domain/exception"
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
