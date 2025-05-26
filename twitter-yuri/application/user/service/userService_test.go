package user

import (
	"Yuri/twitter/core/models"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestCreateUser(t *testing.T) {

	tests := []struct {
		name        string
		inputUser   *models.User
		mockSetup   func(mockRepo *MockUserRepository)
		expectedErr bool
	}{
		{
			name: "User created successfully",
			inputUser: &models.User{
				ID:          "123",
				Name:        "José",
				FollowUsers: []string{},
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("SaveUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedErr: false,
		},
		{
			name:        "User is nil",
			inputUser:   nil,
			mockSetup:   func(mockRepo *MockUserRepository) {},
			expectedErr: true,
		},
		{
			name: "User name is empty",
			inputUser: &models.User{
				ID:          "123",
				Name:        "",
				FollowUsers: []string{},
			},
			mockSetup:   func(mockRepo *MockUserRepository) {},
			expectedErr: true,
		},
		{
			name:      "Repository returns error",
			inputUser: &models.User{ID: "123", Name: "José"},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("SaveUser", mock.AnythingOfType("*models.User")).Return(fmt.Errorf("db error"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			service := NewUserService(mockRepo)
			err := service.CreateUser(tt.inputUser)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFollowUser(t *testing.T) {
	tests := []struct {
		name           string
		idUser         string
		followID       string
		mockUpdateUser func(mockRepo *MockUserRepository)
		expectedError  bool
	}{
		{
			name:     "Success follow user",
			idUser:   "1",
			followID: "2",
			mockUpdateUser: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{},
				}, nil)
				mockRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
					return u.ID == "1" &&
						len(u.FollowUsers) == 1 &&
						u.FollowUsers[0] == "2"
				})).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name:           "Follow himself error",
			idUser:         "1",
			followID:       "1",
			mockUpdateUser: func(mockRepo *MockUserRepository) {},
			expectedError:  true,
		},
		{
			name:     "Already following user",
			idUser:   "1",
			followID: "2",
			mockUpdateUser: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{"2"},
				}, nil)
			},
			expectedError: true,
		},
		{
			name:     "UpdateUser returns error",
			idUser:   "1",
			followID: "2",
			mockUpdateUser: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{},
				}, nil)
				mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(fmt.Errorf("db error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)

			service := &UserServiceImpl{
				repo: mockRepo,
			}

			tt.mockUpdateUser(mockRepo)

			err := service.FollowUser(tt.idUser, tt.followID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUnfollowUser(t *testing.T) {
	tests := []struct {
		name           string
		idUser         string
		unfollowID     string
		mockUpdateUser func(mockRepo *MockUserRepository)
		expectedError  bool
	}{
		{
			name:       "Success unfollow user",
			idUser:     "1",
			unfollowID: "2",
			mockUpdateUser: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{"2"},
				}, nil)
				mockRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
					return u.ID == "1" &&
						len(u.FollowUsers) == 0 &&
						!slices.Contains(u.FollowUsers, "2")
				})).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name:           "Unfollow himself error",
			idUser:         "1",
			unfollowID:     "1",
			mockUpdateUser: func(mockRepo *MockUserRepository) {},
			expectedError:  true,
		},
		{
			name:       "Cannot unfollow user not followed",
			idUser:     "1",
			unfollowID: "2",
			mockUpdateUser: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{"3"},
				}, nil)
			},
			expectedError: true,
		},
		{
			name:       "UpdateUser returns error",
			idUser:     "1",
			unfollowID: "2",
			mockUpdateUser: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{"2"},
				}, nil)
				mockRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(fmt.Errorf("db error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)

			service := &UserServiceImpl{
				repo: mockRepo,
			}

			tt.mockUpdateUser(mockRepo)

			err := service.UnfollowUser(tt.idUser, tt.unfollowID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserById(t *testing.T) {
	tests := []struct {
		name            string
		idUser          string
		mockGetUserById func(mockRepo *MockUserRepository)
		expectedError   bool
	}{
		{
			name:   "Success get user by ID",
			idUser: "1",
			mockGetUserById: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(&models.User{
					ID:          "1",
					Name:        "José",
					FollowUsers: []string{"2"},
				}, nil)
			},
			expectedError: false,
		},
		{
			name:   "Not found user by ID",
			idUser: "2",
			mockGetUserById: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "2").Return(nil, fmt.Errorf("user not found"))
			},
			expectedError: true,
		},

		{
			name:   "Error getting user by ID",
			idUser: "1",
			mockGetUserById: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(nil, fmt.Errorf("database error"))
			},
			expectedError: true,
		},
		{
			name:   "GetUserById returns error",
			idUser: "1",
			mockGetUserById: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetById", "1").Return(nil, fmt.Errorf("db error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)

			service := &UserServiceImpl{
				repo: mockRepo,
			}

			tt.mockGetUserById(mockRepo)

			_, err := service.GetUserById(tt.idUser)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
