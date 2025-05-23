package dtos

type TweetRequest struct {
	UserID  string `json:"user_id" validate:"required,gt=0"`
	Content string `json:"content" validate:"required,max=280"`
}
