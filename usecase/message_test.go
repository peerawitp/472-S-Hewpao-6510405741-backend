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

		// Setup expectations
		mockRepo.On("Store", mock.MatchedBy(func(m *domain.Message) bool {
			return m.UserID == userID && m.ChatID == chatID && m.Content == content
		})).Return(nil).Once()

		// Call the method being tested
		message, err := messageUsecase.CreateMessage(userID, chatID, content)

		// Assert results
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, userID, message.UserID)
		assert.Equal(t, chatID, message.ChatID)
		assert.Equal(t, content, message.Content)

		// Verify expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		userID := "user123"
		chatID := uint(1)
		content := "Hello, world!"
		expectedError := errors.New("database error")

		// Setup expectations
		mockRepo.On("Store", mock.Anything).Return(expectedError).Once()

		// Call the method being tested
		message, err := messageUsecase.CreateMessage(userID, chatID, content)

		// Assert results
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, message)

		// Verify expectations were met
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

		// Setup expectations
		mockRepo.On("GetByChatID", chatID).Return(expectedMessages, nil).Once()

		// Call the method being tested
		messages, err := messageUsecase.GetByChatID(chatID)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, expectedMessages, messages)

		// Verify expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		chatID := "chat123"
		expectedError := errors.New("database error")

		// Setup expectations
		mockRepo.On("GetByChatID", chatID).Return([]domain.Message{}, expectedError).Once()

		// Call the method being tested
		messages, err := messageUsecase.GetByChatID(chatID)

		// Assert results
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, messages)

		// Verify expectations were met
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

		// Setup expectations
		mockRepo.On("GetByID", messageID).Return(expectedMessage, nil).Once()

		// Call the method being tested
		message, err := messageUsecase.GetByID(messageID)

		// Assert results
		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, message)

		// Verify expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error_RepositoryError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		
		messageID := "msg123"
		expectedError := errors.New("database error")

		// Setup expectations
		mockRepo.On("GetByID", messageID).Return((*domain.Message)(nil), expectedError).Once()

		// Call the method being tested
		message, err := messageUsecase.GetByID(messageID)

		// Assert results
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, message)

		// Verify expectations were met
		mockRepo.AssertExpectations(t)
	})
}