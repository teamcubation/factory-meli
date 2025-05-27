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

// MockTweetRepository is a mock implementation of the TweetRepository interface
type MockTweetRepository struct {
	mock.Mock
}

func (m *MockTweetRepository) Save(ctx context.Context, params postgres.TweetSaveParams) (*models.Tweet, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Tweet), args.Error(1)
}

func (m *MockTweetRepository) SaveReply(ctx context.Context, params postgres.TweetSaveReplyParams) (*models.Tweet, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Tweet), args.Error(1)
}

func (m *MockTweetRepository) FindByID(ctx context.Context, params postgres.TweetFindByIDParams) (*models.Tweet, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Tweet), args.Error(1)
}

func (m *MockTweetRepository) FindAllTweetsByUserId(ctx context.Context, params postgres.TweetFindAllTweetsByUserId) ([]*models.Tweet, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Tweet), args.Error(1)
}

func (m *MockTweetRepository) FetchReplies(ctx context.Context, params postgres.TweetFetchRepliesParams) ([]*models.Tweet, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Tweet), args.Error(1)
}

func TestGetAllUserTweets(t *testing.T) {
	tests := []struct {
		name           string
		creatorID      string
		queryParams    string
		mockSetup      func(*MockTweetRepository, *MockUserRepository, uuid.UUID)
		expectedStatus int
		expectedError  bool
		expectedTweets int
	}{
		{
			name:        "Successful retrieval of tweets",
			creatorID:   uuid.New().String(),
			queryParams: "?limit=5&offset=0",
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(&models.User{ID: cID, Email: "test@example.com"}, nil).Once()

				tweetRepo.On("FindAllTweetsByUserId", mock.Anything, postgres.TweetFindAllTweetsByUserId{
					CreatorID: cID,
					Limit:     5,
					Offset:    0,
				}).Return([]*models.Tweet{
					{ID: uuid.New(), Post: "Test tweet 1", CreatorID: cID},
					{ID: uuid.New(), Post: "Test tweet 2", CreatorID: cID},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedTweets: 2,
		},
		{
			name:           "Invalid creator ID",
			creatorID:      "invalid-uuid",
			queryParams:    "",
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected, error before parsing ID
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:        "User not found",
			creatorID:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:        "Database error during user lookup",
			creatorID:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:           "Invalid limit parameter",
			creatorID:      uuid.New().String(),
			queryParams:    "?limit=abc",
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected since limit parsing fails first
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:           "Invalid offset parameter",
			creatorID:      uuid.New().String(),
			queryParams:    "?offset=abc",
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected since offset parsing fails first
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:        "Database error during tweet lookup",
			creatorID:   uuid.New().String(),
			queryParams: "",
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(&models.User{ID: cID, Email: "test@example.com"}, nil).Once()

				tweetRepo.On("FindAllTweetsByUserId", mock.Anything, mock.AnythingOfType("postgres.TweetFindAllTweetsByUserId")).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweets: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Can re-enable now that the flow is clearer

			mockTweetRepo := new(MockTweetRepository)
			mockUserRepo := new(MockUserRepository)

			var parsedCreatorID uuid.UUID
			if tt.creatorID != "invalid-uuid" {
				parsedCreatorID = uuid.MustParse(tt.creatorID)
			}
			tt.mockSetup(mockTweetRepo, mockUserRepo, parsedCreatorID)

			service := NewTweetService(mockTweetRepo, mockUserRepo)

			router := http.NewServeMux()
			var actualTweets []*models.Tweet
			var actualErr error

			router.HandleFunc("/users/{creator_id}/tweets", func(w http.ResponseWriter, r *http.Request) {
				t.Logf("Inside handler for %s: PathValue('creator_id') = %s", tt.name, r.PathValue("creator_id"))
				actualTweets, actualErr = service.GetAllUserTweets(context.Background(), r)
			})

			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.creatorID+"/tweets"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, actualTweets)
			} else {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualTweets)
				assert.Equal(t, tt.expectedTweets, len(actualTweets))
			}

			mockTweetRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestCreateTweet(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		creatorID      string                                                     // Keep as string for path parameter
		mockSetup      func(*MockTweetRepository, *MockUserRepository, uuid.UUID) // Pass parsed UUID for mocks
		expectedStatus int
		expectedError  bool
		expectedTweet  *models.Tweet
	}{
		{
			name:        "Valid tweet creation",
			requestBody: `{"post": "Hello, Twitter!"}`,
			creatorID:   uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(&models.User{ID: cID, Email: "test@example.com"}, nil).Once()

				newTweet := &models.Tweet{
					ID:        uuid.New(),
					Post:      "Hello, Twitter!",
					CreatorID: cID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				tweetRepo.On("Save", mock.Anything, postgres.TweetSaveParams{
					Post:      "Hello, Twitter!",
					CreatorID: cID,
				}).Return(newTweet, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedTweet:  &models.Tweet{},
		},
		{
			name:           "Invalid request payload",
			requestBody:    `{"post": 123}`,                                                                      // Invalid JSON for post
			creatorID:      uuid.New().String(),                                                                  // Valid UUID to allow JSON parsing to fail
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected, JSON parsing fails first
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:           "Empty post content",
			requestBody:    `{"post": ""}`,
			creatorID:      uuid.New().String(),                                                                  // Valid UUID to reach post validation
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected, post validation fails first
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:           "Post content too long",
			requestBody:    `{"post": "` + strings.Repeat("a", 281) + `"}`,                                       // 281 characters
			creatorID:      uuid.New().String(),                                                                  // Valid UUID to reach post validation
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected, post validation fails first
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:           "Invalid creator ID",
			requestBody:    `{"post": "Some tweet"}`,
			creatorID:      "invalid-uuid",
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {}, // No mocks expected, error on parsing ID
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "User not found",
			requestBody: `{"post": "Some tweet"}`,
			creatorID:   uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Database error during user lookup",
			requestBody: `{"post": "Some tweet"}`,
			creatorID:   uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Database error during tweet save",
			requestBody: `{"post": "Some tweet"}`,
			creatorID:   uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, cID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: cID,
				}).Return(&models.User{ID: cID, Email: "test@example.com"}, nil).Once()

				tweetRepo.On("Save", mock.Anything, mock.AnythingOfType("postgres.TweetSaveParams")).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweet:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTweetRepo := new(MockTweetRepository)
			mockUserRepo := new(MockUserRepository)

			var parsedCreatorID uuid.UUID
			if tt.creatorID != "invalid-uuid" {
				parsedCreatorID = uuid.MustParse(tt.creatorID)
			}
			tt.mockSetup(mockTweetRepo, mockUserRepo, parsedCreatorID)

			service := NewTweetService(mockTweetRepo, mockUserRepo)

			router := http.NewServeMux()
			var actualTweet *models.Tweet
			var actualErr error

			router.HandleFunc("/users/{creator_id}/tweets", func(w http.ResponseWriter, r *http.Request) {
				t.Logf("Inside handler for %s: PathValue('creator_id') = %s", tt.name, r.PathValue("creator_id"))
				actualTweet, actualErr = service.CreateTweet(context.Background(), r)
			})

			req := httptest.NewRequest(http.MethodPost, "/users/"+tt.creatorID+"/tweets", strings.NewReader(tt.requestBody))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, actualTweet)
			} else {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualTweet)
			}

			mockTweetRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
