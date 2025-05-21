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

func TestGetUserTimeline(t *testing.T) {
	tests := []struct {
		name           string
		userIDStr      string
		queryParams    string
		mockSetup      func(*MockUserRepository, uuid.UUID)
		expectedStatus int
		expectedError  bool
		expectedResult int
	}{
		{
			name:        "Successful timeline retrieval with defaults",
			userIDStr:   uuid.New().String(),
			queryParams: "", // Uses default limit=10, offset=0, tweet_num=10
			mockSetup: func(repo *MockUserRepository, userID uuid.UUID) {
				// User exists check
				repo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).
					Return(&models.User{ID: userID, Email: "user@example.com"}, nil).Once()

				// GetTimeline call
				expectedTimelineParams := postgres.TimelineParams{
					UserID:            userID,
					TweetsPerFollowed: 10,
					Limit:             10,
					Offset:            0, 
				}
				sampleTimeline := []*models.UserWithTweets{
					{
						User: models.User{ID: uuid.New(), Email: "followed1@example.com"},
						Tweets: []models.Tweet{
							{ID: uuid.New(), Post: "Tweet 1", CreatorID: uuid.New()},
							{ID: uuid.New(), Post: "Tweet 2", CreatorID: uuid.New()},
						},
					},
					{
						User: models.User{ID: uuid.New(), Email: "followed2@example.com"},
						Tweets: []models.Tweet{
							{ID: uuid.New(), Post: "Tweet 3", CreatorID: uuid.New()},
						},
					},
				}
				repo.On("GetTimeline", mock.Anything, expectedTimelineParams).
					Return(sampleTimeline, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: 2, // Number of users with tweets
		},
		{
			name:        "Successful timeline with custom params",
			userIDStr:   uuid.New().String(),
			queryParams: "?limit=5&offset=1&tweet_num=3",
			mockSetup: func(repo *MockUserRepository, userID uuid.UUID) {
				repo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).
					Return(&models.User{ID: userID, Email: "user@example.com"}, nil).Once()

				expectedTimelineParams := postgres.TimelineParams{
					UserID:            userID,
					TweetsPerFollowed: 3,
					Limit:             5,
					Offset:            1,
				}
				sampleTimeline := []*models.UserWithTweets{
					{
						User: models.User{ID: uuid.New(), Email: "followed1@example.com"},
						Tweets: []models.Tweet{
							{ID: uuid.New(), Post: "Custom Tweet A", CreatorID: uuid.New()},
						},
					},
				}
				repo.On("GetTimeline", mock.Anything, expectedTimelineParams).
					Return(sampleTimeline, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: 1,
		},
		{
			name:        "Invalid user ID in path",
			userIDStr:   "invalid-uuid",
			queryParams: "",
			mockSetup:   func(repo *MockUserRepository, userID uuid.UUID) {}, // No mocks expected, fails early on UUID parse
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "User not found",
			userIDStr:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(repo *MockUserRepository, userID uuid.UUID) {
				repo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).
					Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "Database error during user lookup",
			userIDStr:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(repo *MockUserRepository, userID uuid.UUID) {
				repo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).
					Return(nil, sql.ErrConnDone).Once() // Simulate DB error
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "Invalid limit parameter",
			userIDStr:   uuid.New().String(), 
			queryParams: "?limit=abc",
			mockSetup:   func(repo *MockUserRepository, userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "Invalid offset parameter",
			userIDStr:   uuid.New().String(), 
			queryParams: "?offset=xyz",
			mockSetup:   func(repo *MockUserRepository, userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "Invalid tweet_num parameter",
			userIDStr:   uuid.New().String(), 
			queryParams: "?tweet_num=invalid",
			mockSetup:   func(repo *MockUserRepository, userID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "Database error during timeline retrieval",
			userIDStr:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(repo *MockUserRepository, userID uuid.UUID) {
				repo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).
					Return(&models.User{ID: userID, Email: "user@example.com"}, nil).Once()
				repo.On("GetTimeline", mock.Anything, mock.AnythingOfType("postgres.TimelineParams")).
					Return(nil, sql.ErrConnDone).Once() // Simulate DB error during timeline fetch
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedResult: 0,
		},
		{
			name:        "Empty timeline (no followed users or no tweets)",
			userIDStr:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(repo *MockUserRepository, userID uuid.UUID) {
				repo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).
					Return(&models.User{ID: userID, Email: "user@example.com"}, nil).Once()
				repo.On("GetTimeline", mock.Anything, mock.AnythingOfType("postgres.TimelineParams")).
					Return([]*models.UserWithTweets{}, nil).Once() // Return empty slice
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedResult: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(MockUserRepository)

			var parsedUserID uuid.UUID
			if tt.userIDStr != "invalid-uuid" {
				parsedUserID = uuid.MustParse(tt.userIDStr)
			}
			tt.mockSetup(mockRepo, parsedUserID)

			service := NewUserService(mockRepo)

			router := http.NewServeMux()
			var actualTimeline []*models.UserWithTweets
			var actualErr error

			router.HandleFunc("/users/{user_id}/timeline", func(w http.ResponseWriter, r *http.Request) {
				actualTimeline, actualErr = service.GetUserTimeline(context.Background(), r)
			})

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userIDStr+"/timeline"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, actualTimeline)
			} else {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualTimeline)
				assert.Equal(t, tt.expectedResult, len(actualTimeline))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
