package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog" // Import slog
	"net/http"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

// FollowServices defines the interface for follow-related interactions
type FollowServices interface {
	FollowUser(ctx context.Context, r *http.Request) (*models.Follow, error)
	UnfollowUser(ctx context.Context, r *http.Request) error
}

type FollowServiceImpl struct {
	f_repository postgres.FollowRepository
	u_repository postgres.UserRepository
}

func NewFollowService(followRepo postgres.FollowRepository, userRepo postgres.UserRepository) FollowServices {
	return FollowServiceImpl{
		f_repository: followRepo,
		u_repository: userRepo,
	}
}

func (s FollowServiceImpl) FollowUser(ctx context.Context, r *http.Request) (*models.Follow, error) {
	slog.InfoContext(ctx, "Calling service to follow an user", "layer", "service")

	// Parse request body
	type parameters struct {
		FollowerIDStr string `json:"follower_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for FollowUser", "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	// Check if user ids are valid
	followerIDStr := params.FollowerIDStr
	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid follower_id in FollowUser request", "follower_id_str", followerIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid follower_id"}
	}

	followedIDStr := r.PathValue("followed_id")
	followedID, err := uuid.Parse(followedIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid followed_id in FollowUser request", "followed_id_str", followedIDStr, "error", err, "layer", "service")
		return nil, ServiceError{http.StatusBadRequest, "invalid followed_id"}
	}

	if followerID == followedID {
		slog.ErrorContext(ctx, "Attempted to follow self", "follower_id", followerID, "layer", "service")
		return nil, ServiceError{http.StatusConflict, "cannot follow yourself"}
	}

	slog.InfoContext(ctx, "Validating follower and followed users", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followerID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Follower not found", "follower_id", followerID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "follower not found"}
		}
		slog.ErrorContext(ctx, "Error finding follower user by ID", "follower_id", followerID, "error", err, "layer", "service")
		return nil, err
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Followed user not found", "followed_id", followedID, "layer", "service")
			return nil, ServiceError{http.StatusNotFound, "followed not found"}
		}
		slog.ErrorContext(ctx, "Error finding followed user by ID", "followed_id", followedID, "error", err, "layer", "service")
		return nil, err
	}

	// Check if the follow relationship already exists
	slog.InfoContext(ctx, "Checking for existing follow relationship", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	existingFollow, err := s.f_repository.FindByIds(ctx, postgres.FollowFindByIdsParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil && err != sql.ErrNoRows {
		slog.ErrorContext(ctx, "Error checking for existing follow relationship", "follower_id", followerID, "followed_id", followedID, "error", err, "layer", "service")
		return nil, err
	}

	if existingFollow != nil {
		slog.ErrorContext(ctx, "Follow relationship already exists", "follower_id", followerID, "followed_id", followedID, "layer", "service")
		return nil, ServiceError{http.StatusConflict, "follow relationship already exists"}
	}

	// Creating a new follow
	slog.InfoContext(ctx, "Creating new follow relationship", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	follow, err := s.f_repository.Save(ctx, postgres.FollowSaveParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error saving new follow relationship", "follower_id", followerID, "followed_id", followedID, "error", err, "layer", "service")
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully followed user", "follower_id", followerID, "followed_id", followedID, "follow_id", follow.ID, "layer", "service")
	return follow, nil
}

func (s FollowServiceImpl) UnfollowUser(ctx context.Context, r *http.Request) error {
	slog.InfoContext(ctx, "Calling service to unfollow an user", "layer", "service")

	// Parse request body
	type parameters struct {
		FollowerIDStr string `json:"follower_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		slog.ErrorContext(ctx, "Invalid request payload for UnfollowUser", "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()

	followerIDStr := params.FollowerIDStr
	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid follower_id in UnfollowUser request", "follower_id_str", followerIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid follower_id"}
	}

	followedIDStr := r.PathValue("followed_id")
	followedID, err := uuid.Parse(followedIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid followed_id in UnfollowUser request", "followed_id_str", followedIDStr, "error", err, "layer", "service")
		return ServiceError{http.StatusBadRequest, "invalid followed_id"}
	}

	if followerID == followedID {
		slog.ErrorContext(ctx, "Attempted to unfollow self", "follower_id", followerID, "layer", "service")
		return ServiceError{http.StatusBadRequest, "cannot unfollow yourself"}
	}

	slog.InfoContext(ctx, "Validating follower and followed users for unfollow", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followerID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Follower not found for unfollow operation", "follower_id", followerID, "layer", "service")
			return ServiceError{http.StatusNotFound, "follower not found"}
		}
		slog.ErrorContext(ctx, "Error finding follower user by ID for unfollow", "follower_id", followerID, "error", err, "layer", "service")
		return err
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Followed user not found for unfollow operation", "followed_id", followedID, "layer", "service")
			return ServiceError{http.StatusNotFound, "followed not found"}
		}
		slog.ErrorContext(ctx, "Error finding followed user by ID for unfollow", "followed_id", followedID, "error", err, "layer", "service")
		return err
	}

	// Check if the follow relationship exists
	slog.InfoContext(ctx, "Checking if follow relationship exists for unfollow", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	_, err = s.f_repository.FindByIds(ctx, postgres.FollowFindByIdsParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Follow relationship does not exist to unfollow", "follower_id", followerID, "followed_id", followedID, "layer", "service")
			return ServiceError{http.StatusBadRequest, "cannot unfollow a user you don't follow"}
		}
		slog.ErrorContext(ctx, "Error checking for existing follow relationship for unfollow", "follower_id", followerID, "followed_id", followedID, "error", err, "layer", "service")
		return err
	}

	// Deleting the follow
	slog.InfoContext(ctx, "Deleting follow relationship", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	err = s.f_repository.Delete(ctx, postgres.FollowDeleteParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting follow relationship", "follower_id", followerID, "followed_id", followedID, "error", err, "layer", "service")
		return err
	}

	slog.InfoContext(ctx, "Successfully unfollowed user", "follower_id", followerID, "followed_id", followedID, "layer", "service")
	return nil
}
