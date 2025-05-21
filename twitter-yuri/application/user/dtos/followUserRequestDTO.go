package dtos

type UserActionRequestDTO struct {
	UserID       string `json:"user_id" validate:"required"`
	TargetUserId string `json:"target_user_id" validate:"required"`
}
