package request

type CreateUserRequestDTO struct {
	Name        string `json:"name" validate:"required"`
	FollowUsers []int  `json:"follow_users"`
	// PerformedTweets []Tweet
}
