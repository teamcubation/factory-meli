package services

import (
	"context"
	"database/sql"
	"fmt"
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

// MockLikeRepository is a mock implementation of the LikeRepository interface
type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) AddLike(ctx context.Context, params postgres.LikeAddParams) (*models.Like, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Like), args.Error(1)
}

func (m *MockLikeRepository) RemoveLike(ctx context.Context, params postgres.LikeRemoveParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockLikeRepository) HasLike(ctx context.Context, params postgres.LikeHasParams) (bool, error) {
	args := m.Called(ctx, params)
	return args.Bool(0), args.Error(1)
}

// MockRetweetRepository is a mock implementation of the RetweetRepository interface
type MockRetweetRepository struct {
	mock.Mock
}

func (m *MockRetweetRepository) AddRetweet(ctx context.Context, params postgres.RetweetAddParams) (*models.Retweet, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Retweet), args.Error(1)
}

func (m *MockRetweetRepository) HasRetweet(ctx context.Context, params postgres.RetweetHasParams) (bool, error) {
	args := m.Called(ctx, params)
	return args.Bool(0), args.Error(1)
}

func TestLike(t *testing.T) {
	tests := []struct {
		name           string
		userIDStr      string
		tweetIDStr     string
		requestBody    string
		mockSetup      func(*MockUserRepository, *MockLikeRepository, *MockTweetRepository, uuid.UUID, uuid.UUID) // Added userID and tweetID
		expectedStatus int
		expectedError  bool
		expectedLike   *models.Like
	}{
		{
			name:        "Successful like",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(false, nil).Once()
				lRepo.On("AddLike", mock.Anything, postgres.LikeAddParams{UserID: userID, TweetID: tweetID}).Return(&models.Like{ID: uuid.New(), UserID: userID, TweetID: tweetID, CreatedAt: time.Now()}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedLike:   &models.Like{}, // Just check for non-nil
		},
		{
			name:        "Invalid request payload",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": 123}`, // Invalid JSON
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			}, // No mocks expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Invalid user ID in body",
			userIDStr:   "invalid-uuid",
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "invalid-uuid"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			}, // No mocks expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Invalid tweet ID in path",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  "invalid-uuid",
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			}, // No mocks expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "User not found",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Database error during user lookup",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Tweet not found",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Database error during tweet lookup",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Already liked",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(true, nil).Once()
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Database error checking if already liked",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(false, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedLike:   nil,
		},
		{
			name:        "Database error adding like",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(false, nil).Once()
				lRepo.On("AddLike", mock.Anything, postgres.LikeAddParams{UserID: userID, TweetID: tweetID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedLike:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUserRepo := new(MockUserRepository)
			mockLikeRepo := new(MockLikeRepository)
			mockTweetRepo := new(MockTweetRepository)
			mockRetweetRepo := new(MockRetweetRepository) // Not directly used in Like, but required by NewTweetInteractor

			var parsedUserID uuid.UUID
			var parsedTweetID uuid.UUID

			// Only parse if the string is expected to be a valid UUID
			if tt.userIDStr != "invalid-uuid" && tt.userIDStr != "SAME_AS_TWEET_CREATOR" {
				parsedUserID = uuid.MustParse(tt.userIDStr)
			}
			if tt.tweetIDStr != "invalid-uuid" {
				parsedTweetID = uuid.MustParse(tt.tweetIDStr)
			}

			// Call mock setup with the parsed UUIDs
			tt.mockSetup(mockUserRepo, mockLikeRepo, mockTweetRepo, parsedUserID, parsedTweetID)

			service := NewTweetInteractor(mockUserRepo, mockLikeRepo, mockTweetRepo, mockRetweetRepo)

			router := http.NewServeMux()
			var actualLike *models.Like
			var actualErr error

			router.HandleFunc("/tweets/{id}/like", func(w http.ResponseWriter, r *http.Request) {
				actualLike, actualErr = service.Like(context.Background(), r)
			})

			reqBody := fmt.Sprintf(tt.requestBody, tt.userIDStr)
			requestURL := "/tweets/" + tt.tweetIDStr + "/like"

			req := httptest.NewRequest(http.MethodPost, requestURL, strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, actualLike)
			} else {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualLike)
				assert.Equal(t, parsedUserID, actualLike.UserID)
				assert.Equal(t, parsedTweetID, actualLike.TweetID)
			}

			mockUserRepo.AssertExpectations(t)
			mockLikeRepo.AssertExpectations(t)
			mockTweetRepo.AssertExpectations(t)
			mockRetweetRepo.AssertExpectations(t) // Ensure no unexpected calls for retweet repo
		})
	}
}

func TestUnlike(t *testing.T) {
	tests := []struct {
		name           string
		userIDStr      string
		tweetIDStr     string
		requestBody    string
		mockSetup      func(*MockUserRepository, *MockLikeRepository, *MockTweetRepository, uuid.UUID, uuid.UUID) // Added userID and tweetID
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Successful unlike",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(true, nil).Once()
				lRepo.On("RemoveLike", mock.Anything, postgres.LikeRemoveParams{UserID: userID, TweetID: tweetID}).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "Invalid request payload",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": 123}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "Invalid user ID in body",
			userIDStr:   "invalid-uuid",
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "invalid-uuid"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "Invalid tweet ID in path",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  "invalid-uuid",
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "User not found",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:        "Database error during user lookup",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:        "Tweet not found",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:        "Database error during tweet lookup",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:        "Tweet not liked by user",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(false, nil).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:        "Database error checking if liked before unliking",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(false, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:        "Database error removing like",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, lRepo *MockLikeRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID}, nil).Once()
				lRepo.On("HasLike", mock.Anything, postgres.LikeHasParams{UserID: userID, TweetID: tweetID}).Return(true, nil).Once()
				lRepo.On("RemoveLike", mock.Anything, postgres.LikeRemoveParams{UserID: userID, TweetID: tweetID}).Return(sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUserRepo := new(MockUserRepository)
			mockLikeRepo := new(MockLikeRepository)
			mockTweetRepo := new(MockTweetRepository)
			mockRetweetRepo := new(MockRetweetRepository)

			var parsedUserID uuid.UUID
			var parsedTweetID uuid.UUID

			if tt.userIDStr != "invalid-uuid" && tt.userIDStr != "SAME_AS_TWEET_CREATOR" {
				parsedUserID = uuid.MustParse(tt.userIDStr)
			}
			if tt.tweetIDStr != "invalid-uuid" {
				parsedTweetID = uuid.MustParse(tt.tweetIDStr)
			}

			tt.mockSetup(mockUserRepo, mockLikeRepo, mockTweetRepo, parsedUserID, parsedTweetID)

			service := NewTweetInteractor(mockUserRepo, mockLikeRepo, mockTweetRepo, mockRetweetRepo)

			router := http.NewServeMux()
			var actualErr error

			router.HandleFunc("/tweets/{id}/unlike", func(w http.ResponseWriter, r *http.Request) {
				actualErr = service.Unlike(context.Background(), r)
			})

			reqBody := fmt.Sprintf(tt.requestBody, tt.userIDStr)
			requestURL := "/tweets/" + tt.tweetIDStr + "/unlike"

			req := httptest.NewRequest(http.MethodDelete, requestURL, strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
			} else {
				assert.NoError(t, actualErr)
			}

			mockUserRepo.AssertExpectations(t)
			mockLikeRepo.AssertExpectations(t)
			mockTweetRepo.AssertExpectations(t)
			mockRetweetRepo.AssertExpectations(t)
		})
	}
}

func TestRetweet(t *testing.T) {
	tests := []struct {
		name            string
		userIDStr       string
		tweetIDStr      string
		requestBody     string
		mockSetup       func(*MockUserRepository, *MockRetweetRepository, *MockTweetRepository, uuid.UUID, uuid.UUID)
		expectedStatus  int
		expectedError   bool
		expectedRetweet *models.Retweet
	}{
		{
			name:        "Successful retweet",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				// Mock user exists
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				// Mock original tweet exists
				// Ensure CreatorID is different from userID for a valid retweet
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID, CreatorID: uuid.New()}, nil).Once()
				// Mock not already retweeted
				rRepo.On("HasRetweet", mock.Anything, postgres.RetweetHasParams{UserID: userID, TweetID: tweetID}).Return(false, nil).Once()
				// Mock add retweet successful
				rRepo.On("AddRetweet", mock.Anything, postgres.RetweetAddParams{UserID: userID, TweetID: tweetID}).Return(&models.Retweet{ID: uuid.New(), UserID: userID, TweetID: tweetID, CreatedAt: time.Now()}, nil).Once()
			},
			expectedStatus:  http.StatusOK,
			expectedError:   false,
			expectedRetweet: &models.Retweet{}, // Just check for non-nil
		},
		{
			name:        "Invalid request payload",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": 123}`, // Invalid JSON
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			}, // No mocks expected
			expectedStatus:  http.StatusBadRequest,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Invalid user ID in body",
			userIDStr:   "invalid-uuid",
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "invalid-uuid"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			}, // No mocks expected
			expectedStatus:  http.StatusBadRequest,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Invalid tweet ID in path",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  "invalid-uuid",
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
			}, // No mocks expected
			expectedStatus:  http.StatusBadRequest,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "User not found",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus:  http.StatusNotFound,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Database error during user lookup",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Original tweet not found",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus:  http.StatusNotFound,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Database error during original tweet lookup",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Cannot retweet your own tweet",
			userIDStr:   "SAME_AS_TWEET_CREATOR", // Special tag
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				// The tweet's creator ID will be set to userID in the mock
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID, CreatorID: userID}, nil).Once()
			},
			expectedStatus:  http.StatusBadRequest,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Already retweeted",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID, CreatorID: uuid.New()}, nil).Once()
				rRepo.On("HasRetweet", mock.Anything, postgres.RetweetHasParams{UserID: userID, TweetID: tweetID}).Return(true, nil).Once()
			},
			expectedStatus:  http.StatusConflict,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Database error checking if already retweeted",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID, CreatorID: uuid.New()}, nil).Once()
				rRepo.On("HasRetweet", mock.Anything, postgres.RetweetHasParams{UserID: userID, TweetID: tweetID}).Return(false, sql.ErrConnDone).Once()
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   true,
			expectedRetweet: nil,
		},
		{
			name:        "Database error adding retweet",
			userIDStr:   uuid.New().String(),
			tweetIDStr:  uuid.New().String(),
			requestBody: `{"user_id": "%s"}`,
			mockSetup: func(uRepo *MockUserRepository, rRepo *MockRetweetRepository, tRepo *MockTweetRepository, userID, tweetID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: userID}).Return(&models.User{ID: userID}, nil).Once()
				tRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{ID: tweetID}).Return(&models.Tweet{ID: tweetID, CreatorID: uuid.New()}, nil).Once()
				rRepo.On("HasRetweet", mock.Anything, postgres.RetweetHasParams{UserID: userID, TweetID: tweetID}).Return(false, nil).Once()
				rRepo.On("AddRetweet", mock.Anything, postgres.RetweetAddParams{UserID: userID, TweetID: tweetID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedError:   true,
			expectedRetweet: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUserRepo := new(MockUserRepository)
			mockLikeRepo := new(MockLikeRepository) // Not directly used in Retweet, but required by NewTweetInteractor
			mockTweetRepo := new(MockTweetRepository)
			mockRetweetRepo := new(MockRetweetRepository)

			var parsedUserID uuid.UUID
			var parsedTweetID uuid.UUID

			// Parse valid UUIDs for mocks
			if tt.userIDStr != "invalid-uuid" && tt.userIDStr != "SAME_AS_TWEET_CREATOR" {
				parsedUserID = uuid.MustParse(tt.userIDStr)
			}
			if tt.tweetIDStr != "invalid-uuid" {
				parsedTweetID = uuid.MustParse(tt.tweetIDStr)
			}
			if tt.userIDStr == "SAME_AS_TWEET_CREATOR" {
				parsedUserID = uuid.New() // Generate a new one for this test
				// The tweet's creator will be set to this parsedUserID in mockSetup
			}

			tt.mockSetup(mockUserRepo, mockRetweetRepo, mockTweetRepo, parsedUserID, parsedTweetID)

			service := NewTweetInteractor(mockUserRepo, mockLikeRepo, mockTweetRepo, mockRetweetRepo)

			router := http.NewServeMux()
			var actualRetweet *models.Retweet
			var actualErr error

			router.HandleFunc("/tweets/{id}/retweet", func(w http.ResponseWriter, r *http.Request) {
				actualRetweet, actualErr = service.Retweet(context.Background(), r)
			})

			// Use parsedUserID.String() for the request body when it's a valid UUID
			reqBody := fmt.Sprintf(tt.requestBody, parsedUserID.String())
			requestURL := "/tweets/" + parsedTweetID.String() + "/retweet"

			// Handle invalid UUID strings for request URL
			if tt.tweetIDStr == "invalid-uuid" {
				requestURL = "/tweets/invalid-uuid/retweet"
			}

			req := httptest.NewRequest(http.MethodPost, requestURL, strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, actualRetweet)
			} else {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualRetweet)
				assert.Equal(t, parsedUserID, actualRetweet.UserID)
				assert.Equal(t, parsedTweetID, actualRetweet.TweetID)
			}

			mockUserRepo.AssertExpectations(t)
			mockLikeRepo.AssertExpectations(t)
			mockTweetRepo.AssertExpectations(t)
			mockRetweetRepo.AssertExpectations(t)
		})
	}
}
