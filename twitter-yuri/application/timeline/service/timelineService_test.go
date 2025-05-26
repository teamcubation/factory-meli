package service

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

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) SaveUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetById(id string) (*models.User, error) {
	args := m.Called(id)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
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

func TestGetTimeline(t *testing.T) {
	tests := []struct {
		name          string
		userId        string
		page          int
		tweetsPerPage int
		mockSetup     func(mockRepo *MockTweetRepository, mockUserRepo *MockUserRepository)
		expectedErr   bool
	}{
		{
			name:          "Get timeline for user ID successfully",
			userId:        "1",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockTweetRepo *MockTweetRepository, mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					FollowUsers: []string{"2", "3"},
				}, nil)

				mockTweetRepo.On("FindAllByUserTimeline", []string{"2", "3"}, 1, 10).Return([]*models.Tweet{
					{
						ID:          "001",
						Title:       "Tweet Teste",
						Description: "tweet teste description",
						AuthorId:    "2",
						CreatedAt:   time.Now(),
					},
				}, nil)
			},
			expectedErr: false,
		},
		{
			name:          "User not found (nil user)",
			userId:        "2",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockTweetRepo *MockTweetRepository, mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetById", "2").Return(nil, nil)
			},
			expectedErr: true,
		},
		{
			name:          "Error retrieving user",
			userId:        "3",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockTweetRepo *MockTweetRepository, mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetById", "3").Return(nil, fmt.Errorf("db error"))
			},
			expectedErr: true,
		},
		{
			name:          "Error retrieving tweets from timeline",
			userId:        "4",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockTweetRepo *MockTweetRepository, mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetById", "4").Return(&models.User{
					ID:          "4",
					FollowUsers: []string{"5"},
				}, nil)

				mockTweetRepo.On("FindAllByUserTimeline", []string{"5"}, 1, 10).Return(nil, fmt.Errorf("db error"))
			},
			expectedErr: true,
		},
		{
			name:          "No tweets found but no error",
			userId:        "5",
			page:          1,
			tweetsPerPage: 10,
			mockSetup: func(mockTweetRepo *MockTweetRepository, mockUserRepo *MockUserRepository) {
				mockUserRepo.On("GetById", "5").Return(&models.User{
					ID:          "5",
					FollowUsers: []string{"6"},
				}, nil)

				mockTweetRepo.On("FindAllByUserTimeline", []string{"6"}, 1, 10).Return([]*models.Tweet{}, nil)
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTweetRepo := new(MockTweetRepository)
			mockUserRepo := new(MockUserRepository)
			tt.mockSetup(mockTweetRepo, mockUserRepo)

			service := NewTimelineService(mockTweetRepo, mockUserRepo)
			_, err := service.GetTimeline(tt.userId, tt.page, tt.tweetsPerPage)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockTweetRepo.AssertExpectations(t)
		})
	}
}
