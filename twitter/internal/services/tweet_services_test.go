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

func TestReplyToTweet(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		tweetID        string
		mockSetup      func(*MockTweetRepository, *MockUserRepository, uuid.UUID, uuid.UUID)
		expectedStatus int
		expectedError  bool
		expectedTweet  *models.Tweet
	}{
		{
			name:        "Valid reply to tweet",
			requestBody: `{"post": "This is a reply!", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: creatorID,
				}).Return(&models.User{ID: creatorID, Email: "replier@example.com"}, nil).Once()

				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: parentTweetID,
				}).Return(&models.Tweet{ID: parentTweetID, Post: "Original tweet", CreatorID: uuid.New()}, nil).Once()

				newReply := &models.Tweet{
					ID:        uuid.New(),
					Post:      "This is a reply!",
					CreatorID: creatorID,
					ParentID:  uuid.NullUUID{UUID: parentTweetID, Valid: true},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				tweetRepo.On("SaveReply", mock.Anything, postgres.TweetSaveReplyParams{
					Post:      "This is a reply!",
					CreatorID: creatorID,
					ParentID:  parentTweetID,
				}).Return(newReply, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedTweet:  &models.Tweet{}, // Content is checked for nil/not nil, specific fields are not asserted here.
		},
		{
			name:        "Invalid request payload",
			requestBody: `{"post": 123, "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Empty post content",
			requestBody: `{"post": "", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Post content too long",
			requestBody: `{"post": "` + strings.Repeat("b", 281) + `", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Invalid parent tweet ID",
			requestBody: `{"post": "Reply", "creator_id": "%s"}`,
			tweetID:     "invalid-uuid",
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Invalid creator ID in body",
			requestBody: `{"post": "Reply", "creator_id": "invalid-uuid"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Parent tweet not found",
			requestBody: `{"post": "Reply", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: parentTweetID,
				}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Database error during parent tweet lookup",
			requestBody: `{"post": "Reply", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: parentTweetID,
				}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Reply creator not found",
			requestBody: `{"post": "Reply", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: parentTweetID,
				}).Return(&models.Tweet{ID: parentTweetID, Post: "Original", CreatorID: uuid.New()}, nil).Once()

				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: creatorID,
				}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Database error during reply creator lookup",
			requestBody: `{"post": "Reply", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: parentTweetID,
				}).Return(&models.Tweet{ID: parentTweetID, Post: "Original", CreatorID: uuid.New()}, nil).Once()

				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: creatorID,
				}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweet:  nil,
		},
		{
			name:        "Database error during reply save",
			requestBody: `{"post": "Reply", "creator_id": "%s"}`,
			tweetID:     uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, parentTweetID, creatorID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: parentTweetID,
				}).Return(&models.Tweet{ID: parentTweetID, Post: "Original", CreatorID: uuid.New()}, nil).Once()

				userRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{
					ID: creatorID,
				}).Return(&models.User{ID: creatorID, Email: "replier@example.com"}, nil).Once()

				tweetRepo.On("SaveReply", mock.Anything, mock.AnythingOfType("postgres.TweetSaveReplyParams")).Return(nil, sql.ErrConnDone).Once()
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

			parentTweetID := uuid.Nil
			if tt.tweetID != "invalid-uuid" {
				parentTweetID = uuid.MustParse(tt.tweetID)
			}

			creatorID := uuid.New() // Generate a creatorID for the request body string format
			if strings.Contains(tt.requestBody, "%s") {
				tt.requestBody = fmt.Sprintf(tt.requestBody, creatorID.String())
			}
			// If the creator_id in the request body is intentionally invalid, don't parse it.
			if strings.Contains(tt.requestBody, `"creator_id": "invalid-uuid"`) {
				creatorID = uuid.Nil // Set to Nil or another indicator if the test expects invalid UUID
			}

			tt.mockSetup(mockTweetRepo, mockUserRepo, parentTweetID, creatorID)

			service := NewTweetService(mockTweetRepo, mockUserRepo)

			router := http.NewServeMux()
			var actualTweet *models.Tweet
			var actualErr error

			router.HandleFunc("/tweets/{tweet_id}/reply", func(w http.ResponseWriter, r *http.Request) {
				actualTweet, actualErr = service.ReplyToTweet(context.Background(), r)
			})

			req := httptest.NewRequest(http.MethodPost, "/tweets/"+tt.tweetID+"/reply", strings.NewReader(tt.requestBody))
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

func TestGetTweetThread(t *testing.T) {
	tests := []struct {
		name           string
		tweetID        string
		mockSetup      func(*MockTweetRepository, *MockUserRepository, uuid.UUID)
		expectedStatus int
		expectedError  bool
		expectedTweets int
	}{
		{
			name:    "Successful retrieval of tweet thread",
			tweetID: uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, rootTweetID uuid.UUID) {
				rootTweet := &models.Tweet{
					ID:        rootTweetID,
					Post:      "Root tweet",
					CreatorID: uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				reply1 := &models.Tweet{
					ID:        uuid.New(),
					Post:      "Reply 1",
					CreatorID: uuid.New(),
					ParentID:  uuid.NullUUID{UUID: rootTweetID, Valid: true},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				reply2 := &models.Tweet{
					ID:        uuid.New(),
					Post:      "Reply 2",
					CreatorID: uuid.New(),
					ParentID:  uuid.NullUUID{UUID: rootTweetID, Valid: true},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: rootTweetID,
				}).Return(rootTweet, nil).Once()

				tweetRepo.On("FetchReplies", mock.Anything, postgres.TweetFetchRepliesParams{
					RootTweetID: rootTweetID,
				}).Return([]*models.Tweet{rootTweet, reply1, reply2}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedTweets: 3,
		},
		{
			name:           "Invalid root tweet ID",
			tweetID:        "invalid-uuid",
			mockSetup:      func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, rootTweetID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:    "Root tweet not found",
			tweetID: uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, rootTweetID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: rootTweetID,
				}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:    "Database error during root tweet lookup",
			tweetID: uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, rootTweetID uuid.UUID) {
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: rootTweetID,
				}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweets: 0,
		},
		{
			name:    "Database error during fetching replies",
			tweetID: uuid.New().String(),
			mockSetup: func(tweetRepo *MockTweetRepository, userRepo *MockUserRepository, rootTweetID uuid.UUID) {
				rootTweet := &models.Tweet{
					ID:        rootTweetID,
					Post:      "Root tweet",
					CreatorID: uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				tweetRepo.On("FindByID", mock.Anything, postgres.TweetFindByIDParams{
					ID: rootTweetID,
				}).Return(rootTweet, nil).Once()

				tweetRepo.On("FetchReplies", mock.Anything, postgres.TweetFetchRepliesParams{
					RootTweetID: rootTweetID,
				}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedTweets: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTweetRepo := new(MockTweetRepository)
			mockUserRepo := new(MockUserRepository) // User repo is not used in GetTweetThread directly, but it's part of the service struct.

			rootTweetID := uuid.Nil
			if tt.tweetID != "invalid-uuid" {
				rootTweetID = uuid.MustParse(tt.tweetID)
			}
			tt.mockSetup(mockTweetRepo, mockUserRepo, rootTweetID)

			service := NewTweetService(mockTweetRepo, mockUserRepo)

			router := http.NewServeMux()
			var actualTweets []*models.Tweet
			var actualErr error

			router.HandleFunc("/tweets/{tweet_id}/thread", func(w http.ResponseWriter, r *http.Request) {
				actualTweets, actualErr = service.GetTweetThread(context.Background(), r)
			})

			req := httptest.NewRequest(http.MethodGet, "/tweets/"+tt.tweetID+"/thread", nil)
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
