package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hewpao/hewpao-backend/domain"
	"github.com/hewpao/hewpao-backend/usecase"
	"github.com/hewpao/hewpao-backend/usecase/test/mock_repos"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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
		now := time.Now()
		existingChat := &domain.Chat{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:    chatName,
			Message: "Initial message",
		}
		
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
		now := time.Now()
		expectedChat := &domain.Chat{
			Model: gorm.Model{
				ID:        chatID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:    "test-chat",
			Message: "Welcome to test chat!",
		}

		mockRepo.On("GetByID", chatID).Return(expectedChat, nil).Once()

		chat, err := chatUsecase.GetByID(chatID)

		assert.NoError(t, err)
		assert.Equal(t, expectedChat, chat)
		assert.Equal(t, chatID, chat.ID)
		assert.Equal(t, "test-chat", chat.Name)
		assert.Equal(t, "Welcome to test chat!", chat.Message)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatID := uint(1)
		expectedError := errors.New("database error")
		now := time.Now()
		expectedChat := &domain.Chat{
			Model: gorm.Model{
				ID:        0,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Message: "",
		}  

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
		now := time.Now()
		expectedChat := &domain.Chat{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Name:    chatName,
			Message: "Welcome to test chat!",
		}

		// The GetByName method is called twice in our usecase implementation
		mockRepo.On("GetByName", chatName).Return(expectedChat).Times(2)

		chat := chatUsecase.GetByName(chatName)

		assert.NotNil(t, chat)
		assert.Equal(t, expectedChat, chat)
		assert.Equal(t, uint(1), chat.ID)
		assert.Equal(t, chatName, chat.Name)
		assert.Equal(t, "Welcome to test chat!", chat.Message)

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