package dtos

import "github.com/google/uuid"

type CreateUserRequestDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type FollowRequestDTO struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

type UnfollowRequestDTO struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

type TimelineRequestDTO struct {
	UserID string `json:"user_id"`
}
