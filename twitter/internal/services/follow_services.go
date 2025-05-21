package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

type FollowServiceImpl struct {
	f_repository postgres.FollowRepository
	u_repository postgres.UserRepository
}

func NewFollowService(followRepo postgres.FollowRepository, userRepo postgres.UserRepository) FollowServiceImpl {
	return FollowServiceImpl{
		f_repository: followRepo,
		u_repository: userRepo,
	}
}

func (s *FollowServiceImpl) FollowUser(ctx context.Context, r *http.Request) (*models.Follow, error) {
	// Parse request body
	type parameters struct {
		FollowerIDStr string `json:"follower_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()
	
	// Check if user ids are valid
	followerIDStr := params.FollowerIDStr
	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid follower_id"}
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followerID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ServiceError{http.StatusNotFound, "follower not found"}
		}
		return nil, err
	}
	
	followedIDStr := r.PathValue("followed_id")
	followedID, err := uuid.Parse(followedIDStr)
	if err != nil {
		return nil, ServiceError{http.StatusBadRequest, "invalid followed_id"}
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ServiceError{http.StatusNotFound, "followed not found"}
		}
		return nil, err
	}
	
	if followerID == followedID {
		return nil, ServiceError{http.StatusConflict, "cannot follow yourself"}
	}
	
	// Check if the follow relationship already exists
	existingFollow, err := s.f_repository.FindByIds(ctx, postgres.FollowFindByIdsParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingFollow != nil {
		return nil, ServiceError{http.StatusConflict, "follow relationship already exists"}
	}

	// Creating a new follow
	follow, err := s.f_repository.Save(ctx, postgres.FollowSaveParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (s *FollowServiceImpl) UnfollowUser(ctx context.Context, r *http.Request) error {
	// Parse request body
	type parameters struct {
		FollowerIDStr string `json:"follower_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return ServiceError{http.StatusBadRequest, "invalid request payload"}
	}
	defer r.Body.Close()
	
	followerIDStr := params.FollowerIDStr
	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		return ServiceError{http.StatusBadRequest, "invalid follower_id"}
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followerID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
		  return ServiceError{http.StatusNotFound, "follower not found"}
		}
		return err
	}
	
	followedIDStr := r.PathValue("followed_id")
	followedID, err := uuid.Parse(followedIDStr)
	if err != nil {
		return ServiceError{http.StatusBadRequest, "invalid followed_id"}
	}

	_, err = s.u_repository.FindByID(ctx, postgres.UserFindByIDParams{
		ID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
		  return ServiceError{http.StatusNotFound, "followed not found"}
		}
		return err
	}
	
	if followerID == followedID {
		return ServiceError{http.StatusBadRequest, "cannot unfollow yourself"}
	}
	
	// Check if the follow relationship exists
	_, err = s.f_repository.FindByIds(ctx, postgres.FollowFindByIdsParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ServiceError{http.StatusBadRequest, "cannot unfollow a user you don't follow"}
		}
		return err
	}

	// Deleting the follow
	err = s.f_repository.Delete(ctx, postgres.FollowDeleteParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		return err
	}

	return nil
}
