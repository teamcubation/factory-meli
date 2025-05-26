package tweet

import (
	"Yuri/twitter/core/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTweetRepository struct {
	mock.Mock
}

func (m *MockTweetRepository) SaveTweet(user *models.Tweet) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockTweetRepository) FindAllByUser(userID string, page, tweetsPerPage int) ([]*models.Tweet, error) {
	args := m.Called(userID, page, tweetsPerPage)
	if users := args.Get(0); users != nil {
		return users.([]*models.Tweet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTweetRepository) FindAllByUserTimeline(userIds []string, page int, tweetsPerPage int) ([]*models.Tweet, error) {
	args := m.Called(userIds, page, tweetsPerPage)
	if users := args.Get(0); users != nil {
		return users.([]*models.Tweet), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestCreateTweet(t *testing.T) {
	tests := []struct {
		name        string
		tweetTest   *models.Tweet
		mockSetup   func(mockRepo *MockTweetRepository)
		expectedErr bool
	}{
		{
			name: "Tweet created successfully",
			tweetTest: &models.Tweet{
				ID:          "001",
				Title:       "Test Tweet",
				Description: "This is a test tweet description.",
				AuthorId:    "1",
				CreatedAt:   time.Now(),
			},
			mockSetup: func(mockRepo *MockTweetRepository) {
				mockRepo.On("SaveTweet", mock.AnythingOfType("*models.Tweet")).Return(nil)
			},
			expectedErr: false,
		},
		{
			name:        "Tweet is nil",
			tweetTest:   nil,
			mockSetup:   func(mockRepo *MockTweetRepository) {},
			expectedErr: true,
		},
		{
			name: "Repository returns error",
			tweetTest: &models.Tweet{
				ID:          "001",
				Title:       "Test Tweet",
				Description: "This is a test tweet description.",
				AuthorId:    "1",
				CreatedAt:   time.Now(),
			},
			mockSetup: func(mockRepo *MockTweetRepository) {
				mockRepo.On("SaveTweet", mock.AnythingOfType("*models.Tweet")).Return(fmt.Errorf("db error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTweetRepository)
			tt.mockSetup(mockRepo)

			service := NewTweetService(mockRepo)
			err := service.CreateTweet(tt.tweetTest)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListTweetsByUserID(t *testing.T) {
	tests := []struct {
		name          string
		userId        string
		page          int
		tweetsPerPage int
		mockSetup     func(mockRepo *MockTweetRepository)
		expectedErr   bool
	}{
		{
			name:          "List tweets by user ID successfully",
			userId:        "1",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockRepo *MockTweetRepository) {
				mockRepo.On("FindAllByUser", "1", 1, 10).Return([]*models.Tweet{
					{
						ID:          "001",
						Title:       "Test Tweet",
						Description: "This is a test tweet description.",
						AuthorId:    "1",
						CreatedAt:   time.Now(),
					},
				}, nil)
			},
			expectedErr: false,
		},
		{
			name:          "No tweets found for user ID",
			userId:        "2",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockRepo *MockTweetRepository) {
				mockRepo.On("FindAllByUser", "2", 1, 10).Return([]*models.Tweet{}, nil)
			},
			expectedErr: true,
		},
		{
			name:          "Repository returns error",
			userId:        "2",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockRepo *MockTweetRepository) {
				mockRepo.On("FindAllByUser", "2", 1, 10).Return(nil, fmt.Errorf("db error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTweetRepository)
			tt.mockSetup(mockRepo)

			service := NewTweetService(mockRepo)
			_, err := service.ListTweetsByUserID(tt.userId, tt.page, tt.tweetsPerPage)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
