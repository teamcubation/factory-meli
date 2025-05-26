package services

import (
	"context"
	"database/sql"
	"errors"
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

// MockFollowRepository is a mock implementation of the FollowRepository interface
type MockFollowRepository struct {
	mock.Mock
}

func (m *MockFollowRepository) Save(ctx context.Context, params postgres.FollowSaveParams) (*models.Follow, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Follow), args.Error(1)
}

func (m *MockFollowRepository) Delete(ctx context.Context, params postgres.FollowDeleteParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockFollowRepository) FindByIds(ctx context.Context, params postgres.FollowFindByIdsParams) (*models.Follow, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Follow), args.Error(1)
}

func (m *MockFollowRepository) FindFollowing(ctx context.Context, params postgres.FollowFindFollowingParams) ([]*models.User, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.User), args.Error(1)
}

func TestFollowUser(t *testing.T) {
	tests := []struct {
		name           string
		followerIDStr  string
		followedIDStr  string
		requestBody    string
		mockSetup      func(*MockFollowRepository, *MockUserRepository, uuid.UUID, uuid.UUID)
		expectedStatus int
		expectedError  bool
		expectedFollow *models.Follow
	}{
		{
			name:          "Successful follow",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`, // Will be formatted with followerIDStr
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				// Mock follower exists
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				// Mock followed exists
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				// Mock follow relationship does not exist
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(nil, sql.ErrNoRows).Once()
				// Mock save successful
				fRepo.On("Save", mock.Anything, postgres.FollowSaveParams{FollowerID: followerID, FollowedID: followedID}).Return(&models.Follow{ID: uuid.New(), FollowerID: followerID, FollowedID: followedID, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			expectedFollow: &models.Follow{}, // Just check for non-nil
		},
		{
			name:           "Invalid request payload",
			followerIDStr:  uuid.New().String(),
			followedIDStr:  uuid.New().String(),
			requestBody:    `{"follower_id": 123}`,                                                                            // Invalid JSON
			mockSetup:      func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {}, // No mocks expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:           "Invalid follower ID in body",
			followerIDStr:  "invalid-uuid",
			followedIDStr:  uuid.New().String(), // Valid followed ID to reach follower parsing
			requestBody:    `{"follower_id": "invalid-uuid"}`,
			mockSetup:      func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {}, // No mocks expected
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Follower not found",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Database error during follower lookup",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Invalid followed ID in path",
			followerIDStr: uuid.New().String(), // Valid follower ID to allow parsing to proceed
			followedIDStr: "invalid-uuid",
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				// NO mocks expected. The service will return 'invalid followed_id'
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Followed not found",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Database error during followed lookup",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Cannot follow yourself",
			followerIDStr: uuid.New().String(),
			followedIDStr: "SAME_AS_FOLLOWER", // Special tag for this test case
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				// IMPORTANT: No FindByID calls are expected here as the service checks for equality
				// IMMEDIATELY after parsing both IDs, before any DB checks.
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Follow relationship already exists",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(&models.Follow{ID: uuid.New()}, nil).Once() // Found existing
			},
			expectedStatus: http.StatusConflict,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Database error during FindByIds check",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(nil, sql.ErrConnDone).Once() // Database error
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedFollow: nil,
		},
		{
			name:          "Database error during Save",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(nil, sql.ErrNoRows).Once()
				fRepo.On("Save", mock.Anything, postgres.FollowSaveParams{FollowerID: followerID, FollowedID: followedID}).Return(nil, sql.ErrConnDone).Once() // Save fails
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
			expectedFollow: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockFollowRepo := new(MockFollowRepository)
			mockUserRepo := new(MockUserRepository)

			var parsedFollowerID uuid.UUID
			var parsedFollowedID uuid.UUID

			// Handle special case for "Cannot follow/unfollow yourself"
			if tt.followedIDStr == "SAME_AS_FOLLOWER" {
				// For this specific test, both IDs will be the same
				parsedFollowerID = uuid.MustParse(tt.followerIDStr)
				parsedFollowedID = parsedFollowerID
			} else {
				// Parse valid UUIDs for mocks
				if tt.followerIDStr != "invalid-uuid" {
					parsedFollowerID = uuid.MustParse(tt.followerIDStr)
				}
				// Only parse followedID if it's a valid UUID string and not the special tag
				if tt.followedIDStr != "invalid-uuid" && tt.followedIDStr != "SAME_AS_FOLLOWER" {
					parsedFollowedID = uuid.MustParse(tt.followedIDStr)
				}
			}

			tt.mockSetup(mockFollowRepo, mockUserRepo, parsedFollowerID, parsedFollowedID)

			service := NewFollowService(mockFollowRepo, mockUserRepo)

			router := http.NewServeMux()
			var actualFollow *models.Follow
			var actualErr error

			router.HandleFunc("/follows/{followed_id}", func(w http.ResponseWriter, r *http.Request) {
				actualFollow, actualErr = service.FollowUser(context.Background(), r)
			})

			// Format request body if it contains a placeholder
			reqBody := tt.requestBody
			if strings.Contains(reqBody, "%s") {
				reqBody = fmt.Sprintf(reqBody, tt.followerIDStr)
			}

			// Ensure the request URL for "SAME_AS_FOLLOWER" case uses the actual parsed ID
			requestURL := "/follows/" + tt.followedIDStr
			if tt.followedIDStr == "SAME_AS_FOLLOWER" {
				requestURL = "/follows/" + parsedFollowedID.String()
			}

			req := httptest.NewRequest(http.MethodPost, requestURL, strings.NewReader(reqBody))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if tt.expectedError {
				assert.Error(t, actualErr)
				if svcErr, ok := actualErr.(ServiceError); ok {
					assert.Equal(t, tt.expectedStatus, svcErr.StatusCode)
				}
				assert.Nil(t, actualFollow)
			} else {
				assert.NoError(t, actualErr)
				assert.NotNil(t, actualFollow)
				assert.Equal(t, parsedFollowerID, actualFollow.FollowerID)
				assert.Equal(t, parsedFollowedID, actualFollow.FollowedID)
			}

			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUnfollowUser(t *testing.T) {
	tests := []struct {
		name           string
		followerIDStr  string
		followedIDStr  string
		requestBody    string
		mockSetup      func(*MockFollowRepository, *MockUserRepository, uuid.UUID, uuid.UUID)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:          "Successful unfollow",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				// Mock follower exists
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				// Mock followed exists
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				// Mock follow relationship exists
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(&models.Follow{ID: uuid.New()}, nil).Once()
				// Mock delete successful
				fRepo.On("Delete", mock.Anything, postgres.FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}).Return(nil).Once()
			},
			expectedStatus: http.StatusOK, // Unfollow typically returns 200 OK with no body or 204 No Content
			expectedError:  false,
		},
		{
			name:           "Invalid request payload",
			followerIDStr:  uuid.New().String(),
			followedIDStr:  uuid.New().String(),
			requestBody:    `{"follower_id": 123}`, // Invalid JSON
			mockSetup:      func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Invalid follower ID in body",
			followerIDStr:  "invalid-uuid",
			followedIDStr:  uuid.New().String(),
			requestBody:    `{"follower_id": "invalid-uuid"}`,
			mockSetup:      func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:          "Follower not found",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:          "Database error during follower lookup",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:           "Invalid followed ID in path",
			followerIDStr:  uuid.New().String(),
			followedIDStr:  "invalid-uuid",
			requestBody:    `{"follower_id": "%s"}`,
			mockSetup:      func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:          "Followed not found",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:          "Database error during followed lookup",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:           "Cannot unfollow yourself",
			followerIDStr:  uuid.New().String(),
			followedIDStr:  "SAME_AS_FOLLOWER",
			requestBody:    `{"follower_id": "%s"}`,
			mockSetup:      func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:          "Cannot unfollow a user you don't follow",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(nil, sql.ErrNoRows).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:          "Database error during FindByIds check",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(nil, sql.ErrConnDone).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:          "Database error during Delete",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(&models.Follow{ID: uuid.New()}, nil).Once()
				fRepo.On("Delete", mock.Anything, postgres.FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}).Return(sql.ErrConnDone).Once() // Delete fails
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:          "Follow relationship not found by Delete (return error)",
			followerIDStr: uuid.New().String(),
			followedIDStr: uuid.New().String(),
			requestBody:   `{"follower_id": "%s"}`,
			mockSetup: func(fRepo *MockFollowRepository, uRepo *MockUserRepository, followerID, followedID uuid.UUID) {
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followerID}).Return(&models.User{ID: followerID}, nil).Once()
				uRepo.On("FindByID", mock.Anything, postgres.UserFindByIDParams{ID: followedID}).Return(&models.User{ID: followedID}, nil).Once()
				fRepo.On("FindByIds", mock.Anything, postgres.FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}).Return(&models.Follow{ID: uuid.New()}, nil).Once()
				// This mocks the specific error returned by your repo if RowsAffected is 0
				fRepo.On("Delete", mock.Anything, postgres.FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}).Return(errors.New("follow relationship not found")).Once()
			},
			expectedStatus: http.StatusBadRequest, // Your service maps this specific error to 400
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockFollowRepo := new(MockFollowRepository)
			mockUserRepo := new(MockUserRepository)

			var parsedFollowerID uuid.UUID
			var parsedFollowedID uuid.UUID

			// Handle special case for "Cannot follow/unfollow yourself"
			if tt.followedIDStr == "SAME_AS_FOLLOWER" {
				parsedFollowerID = uuid.MustParse(tt.followerIDStr)
				parsedFollowedID = parsedFollowerID // Set followedID to be the same as followerID
			} else {
				// Parse valid UUIDs for mocks
				if tt.followerIDStr != "invalid-uuid" {
					parsedFollowerID = uuid.MustParse(tt.followerIDStr)
				}
				if tt.followedIDStr != "invalid-uuid" && tt.followedIDStr != "SAME_AS_FOLLOWER" {
					parsedFollowedID = uuid.MustParse(tt.followedIDStr)
				}
			}

			tt.mockSetup(mockFollowRepo, mockUserRepo, parsedFollowerID, parsedFollowedID)

			service := NewFollowService(mockFollowRepo, mockUserRepo)

			router := http.NewServeMux()
			var actualErr error // Unfollow returns an error, not a model

			router.HandleFunc("/follows/{followed_id}", func(w http.ResponseWriter, r *http.Request) {
				actualErr = service.UnfollowUser(context.Background(), r)
			})

			// Format request body if it contains a placeholder
			reqBody := tt.requestBody
			if strings.Contains(reqBody, "%s") {
				reqBody = fmt.Sprintf(reqBody, tt.followerIDStr)
			}

			// Ensure the request URL for "SAME_AS_FOLLOWER" case uses the actual parsed ID
			requestURL := "/follows/" + tt.followedIDStr
			if tt.followedIDStr == "SAME_AS_FOLLOWER" {
				requestURL = "/follows/" + parsedFollowedID.String() // Use the actual UUID string
			}

			req := httptest.NewRequest(http.MethodDelete, requestURL, strings.NewReader(reqBody)) // Unfollow uses DELETE method
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

			mockFollowRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
