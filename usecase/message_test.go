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

func TestCreateMessage(t *testing.T) {
	mockRepo := new(mock_repos.MockMessageRepository)
	messageUsecase := usecase.NewMessageService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		userID := "1"
		chatID := uint(1)
		content := "Hello, world!"

		mockRepo.On("Store", mock.MatchedBy(func(m *domain.Message) bool {
			return m.UserID == userID && m.ChatID == chatID && m.Content == content
		})).Return(nil).Once()

		message, err := messageUsecase.CreateMessage(userID, chatID, content)

		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, userID, message.UserID)
		assert.Equal(t, chatID, message.ChatID)
		assert.Equal(t, content, message.Content)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		userID := "user123"
		chatID := uint(1)
		content := "Hello, world!"
		expectedError := errors.New("database error")

		mockRepo.On("Store", mock.Anything).Return(expectedError).Once()

		
		message, err := messageUsecase.CreateMessage(userID, chatID, content)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, message)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetByChatID(t *testing.T) {
	mockRepo := new(mock_repos.MockMessageRepository)
	messageUsecase := usecase.NewMessageService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatID := "chat123"
		expectedMessages := []domain.Message{
			{UserID: "user1", ChatID: 1, Content: "Message 1"},
			{UserID: "user2", ChatID: 1, Content: "Message 2"},
		}

		
		mockRepo.On("GetByChatID", chatID).Return(expectedMessages, nil).Once()

		
		messages, err := messageUsecase.GetByChatID(chatID)

		
		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, messages)

		
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatID := "chat123"
		expectedError := errors.New("database error")

		
		mockRepo.On("GetByChatID", chatID).Return([]domain.Message{}, expectedError).Once()

		
		messages, err := messageUsecase.GetByChatID(chatID)

		
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, messages)

		
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockRepo := new(mock_repos.MockMessageRepository)
	messageUsecase := usecase.NewMessageService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		messageID := "msg123"
		expectedMessage := &domain.Message{
			UserID:  "user1",
			ChatID:  1,
			Content: "Hello, world!",
		}

		
		mockRepo.On("GetByID", messageID).Return(expectedMessage, nil).Once()

		
		message, err := messageUsecase.GetByID(messageID)

		
		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, message)

		
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		messageID := "msg123"
		expectedError := errors.New("database error")

		
		mockRepo.On("GetByID", messageID).Return((*domain.Message)(nil), expectedError).Once()

		
		message, err := messageUsecase.GetByID(messageID)

		
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, message)

		
		mockRepo.AssertExpectations(t)
	})
}