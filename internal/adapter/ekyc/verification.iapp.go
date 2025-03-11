package ekyc

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/hewpao/hewpao-backend/config"
	"github.com/hewpao/hewpao-backend/dto"
	"github.com/hewpao/hewpao-backend/repository"
)

type iappVerificationRepo struct {
	cfg     *config.Config
	httpCli *http.Client
}

func NewIappVerificationRepo(cfg *config.Config, httpCli *http.Client) repository.EKYCRepository {
	return &iappVerificationRepo{
		cfg:     cfg,
		httpCli: httpCli,
	}
}

func (i *iappVerificationRepo) Verify(file *multipart.FileHeader) (*dto.EKYCResponseDTO, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, err
	}

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	_, err = io.Copy(part, fileReader)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", i.cfg.KYCApiUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("apikey", i.cfg.KYCApiKey)

	res, err := i.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result dto.EKYCResponseDTO
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
