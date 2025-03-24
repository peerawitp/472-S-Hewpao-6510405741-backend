package usecase_test

import (
	"errors"
	"testing"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateChat(t *testing.T) {
	mockRepo := new(mock_repos.MockChatRepository)
	chatUsecase := usecase.NewChatService(mockRepo)

	t.Run("Success_NewChat", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatName := "test-chat"
		
		mockRepo.On("GetByName", chatName).Return(nil).Once()
		mockRepo.On("Create", mock.MatchedBy(func(c *domain.Chat) bool {
			return c.Name == chatName
		})).Return(nil).Once()

		err := chatUsecase.CreateChat(chatName)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Success_ExistingChat", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatName := "existing-chat"
		existingChat := &domain.Chat{Name: chatName}
		
		mockRepo.On("GetByName", chatName).Return(existingChat).Once()

		err := chatUsecase.CreateChat(chatName)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatName := "test-chat"
		expectedError := errors.New("database error")
		
		mockRepo.On("GetByName", chatName).Return(nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*domain.Chat")).Return(expectedError).Once()

		err := chatUsecase.CreateChat(chatName)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetChatByID(t *testing.T) {
	mockRepo := new(mock_repos.MockChatRepository)
	chatUsecase := usecase.NewChatService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatID := uint(1)
		expectedChat := &domain.Chat{
			Name: "test-chat",
		}

		mockRepo.On("GetByID", chatID).Return(expectedChat, nil).Once()

		chat, err := chatUsecase.GetByID(chatID)

		assert.NoError(t, err)
		assert.Equal(t, expectedChat, chat)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatID := uint(1)
		expectedError := errors.New("database error")
		expectedChat := &domain.Chat{}  

		mockRepo.On("GetByID", chatID).Return(expectedChat, expectedError).Once()

		_, err := chatUsecase.GetByID(chatID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetByName(t *testing.T) {
	mockRepo := new(mock_repos.MockChatRepository)
	chatUsecase := usecase.NewChatService(mockRepo)

	t.Run("Success_ChatExists", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatName := "test-chat"
		expectedChat := &domain.Chat{
			Name: chatName,
		}

		mockRepo.On("GetByName", chatName).Return(expectedChat).Times(2)

		chat := chatUsecase.GetByName(chatName)

		assert.NotNil(t, chat)
		assert.Equal(t, expectedChat, chat)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_ChatNotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatName := "nonexistent-chat"

		mockRepo.On("GetByName", chatName).Return(nil).Once()

		chat := chatUsecase.GetByName(chatName)

		assert.Nil(t, chat)

		mockRepo.AssertExpectations(t)
	})
}