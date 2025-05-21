package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/core/ports/output/postgres"
)

type FollowServiceImpl struct {
	repository postgres.FollowRepository
}

func NewFollowService(repo postgres.FollowRepository) FollowServiceImpl {
	return FollowServiceImpl{
		repository: repo,
	}
}

func (s *FollowServiceImpl) FollowUser(ctx context.Context, followerIDStr, followedIDStr string) (*models.Follow, error) {
	// First check if the follow relationship already exists
	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		return nil, errors.New("invalid follower id")
	}

	followedID, err := uuid.Parse(followedIDStr)
	if err != nil {
		return nil, errors.New("invalid followed id")
	}

	if followerID == followedID {
		return nil, errors.New("cannot follow yourself")
	}

	existingFollow, err := s.repository.FindByIds(ctx, postgres.FollowFindByIdsParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingFollow != nil {
		return nil, errors.New("follow relationship already exists")
	}

	// Creating a new follow
	follow, err := s.repository.Save(ctx, postgres.FollowSaveParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		return nil, err
	}

	return follow, nil
}

func (s *FollowServiceImpl) UnfollowUser(ctx context.Context, followerIDStr, followedIDStr string) (error) {
	// First check if the follow relationship already exists
	followerID, err := uuid.Parse(followerIDStr)
	if err != nil {
		return errors.New("invalid follower id")
	}

	followedID, err := uuid.Parse(followedIDStr)
	if err != nil {
		return errors.New("invalid followed id")
	}

	if followerID == followedID {
		return errors.New("cannot unfollow yourself")
	}

	_, err = s.repository.FindByIds(ctx, postgres.FollowFindByIdsParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("cannot unfollow a user you don't follow")
		}
		return err
	}

	// Deleting the follow
	err = s.repository.Delete(ctx, postgres.FollowDeleteParams{
		FollowerID: followerID,
		FollowedID: followedID,
	})
	if err != nil {
		return err
	}

	return nil
}
