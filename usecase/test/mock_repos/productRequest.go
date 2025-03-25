package mock_repos

import (
	"github.com/hewpao/hewpao-backend/domain"
	"github.com/stretchr/testify/mock"
)

type MockProductRequestRepo struct {
	mock.Mock
}

func (mp *MockProductRequestRepo) Create(productRequest *domain.ProductRequest) error {
	args := mp.Called(productRequest)
	return args.Error(0)
}

func (mp *MockProductRequestRepo) FindByID(id int) (*domain.ProductRequest, error) {
	args := mp.Called(id)
	ret0, _ := args.Get(0).(*domain.ProductRequest)
	ret1, _ := args.Get(1).(error)
	return ret0, ret1
}

func (mp *MockProductRequestRepo) FindByUserID(id string) ([]domain.ProductRequest, error) {
	args := mp.Called(id)
	return args.Get(0).([]domain.ProductRequest), args.Error(1)
}

func (mp *MockProductRequestRepo) FindByOfferUserID(id string) ([]domain.ProductRequest, error) {
	args := mp.Called(id)
	return args.Get(0).([]domain.ProductRequest), args.Error(1)
}

func (mp *MockProductRequestRepo) FindPaginatedProductRequests(page, limit int) ([]domain.ProductRequest, int64, error) {
	args := mp.Called(page, limit)
	return args.Get(0).([]domain.ProductRequest), args.Get(1).(int64), args.Error(2)
}

func (mp *MockProductRequestRepo) Update(productRequest *domain.ProductRequest) error {
	args := mp.Called(productRequest)
	return args.Error(0)
}

func (mp *MockProductRequestRepo) IsOwnedByUser(prID int, userID string) (bool, error) {
	args := mp.Called(prID, userID)
	return args.Bool(0), args.Error(1)
}

func (mp *MockProductRequestRepo) Delete(productRequest *domain.ProductRequest) error {
	args := mp.Called(productRequest)
	return args.Error(0)
}
