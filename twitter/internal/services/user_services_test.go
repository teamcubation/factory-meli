package services

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, params postgres.UserFindByEmailParams) (*models.User, error) {
	args := m.Called(ctx, params)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, params postgres.UserFindByIDParams) (*models.User, error) {
	args := m.Called(ctx, params)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, params postgres.UserSaveParams) (*models.User, error) {
	args := m.Called(ctx, params)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetTimeline(ctx context.Context, params postgres.TimelineParams) ([]*models.UserWithTweets, error) {
	args := m.Called(ctx, params)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).([]*models.UserWithTweets), args.Error(1)
}

// Test scenarios
func TestCreateUser(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
		expectedUser   *models.User
	}{
		{
			name:        "Valid user creation",
			requestBody: `{"email": "test@example.com"}`,
			mockSetup: func(repo *MockUserRepository) {
				// Mock FindByEmail to return no user (email not found)
				repo.On("FindByEmail", mock.Anything, postgres.UserFindByEmailParams{
					Email: "test@example.com",
				}).Return(nil, sql.ErrNoRows)
				
				// Mock Save to return a new user
				newUser := &models.User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.On("Save", mock.Anything, postgres.UserSaveParams{
					Email: "test@example.com",
				}).Return(newUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedUser:   &models.User{},  // We'll just check that it's not nil
		},
		{
			name:        "Invalid email format",
			requestBody: `{"email": "invalid-email"}`,
			mockSetup:   func(repo *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedUser:   nil,
		},
		{
			name:        "Email already exists",
			requestBody: `{"email": "existing@example.com"}`,
			mockSetup: func(repo *MockUserRepository) {
				// Mock FindByEmail to return an existing user
				existingUser := &models.User{
					ID:        uuid.New(),
					Email:     "existing@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.On("FindByEmail", mock.Anything, postgres.UserFindByEmailParams{
					Email: "existing@example.com",
				}).Return(existingUser, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			expectedUser:   nil,
		},
		{
			name:        "Empty email",
			requestBody: `{"email": ""}`,
			mockSetup:   func(repo *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedUser:   nil,
		},
		{
			name:        "Database error during lookup",
			requestBody: `{"email": "test@example.com"}`,
			mockSetup: func(repo *MockUserRepository) {
				// Mock FindByEmail to return a database error
				repo.On("FindByEmail", mock.Anything, postgres.UserFindByEmailParams{
					Email: "test@example.com",
				}).Return(nil, sql.ErrConnDone)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedUser:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)
			service := NewUserService(mockRepo)
			
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.requestBody))
			user, err := service.CreateUser(context.Background(), req)
			
			if tt.expectedError {
				assert.Error(t, err)
				if svcErr, ok := err.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}
