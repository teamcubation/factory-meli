package dtos

type FollowRequest struct {
	UserID      string `json:"user_id"`
	FollowingID string `json:"following_id"`
}
