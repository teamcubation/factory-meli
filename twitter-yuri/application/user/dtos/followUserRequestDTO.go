package dtos

// "Yuri/twitter/core/models"

type UserActionRequestDTO struct {
	UserID       string `json:"user_id" validate:"required"`
	TargetUserId string `json:"target_user_id" validate:"required"`
}

// func FollowToModel(dto *FollowUserRequestDTO) *models.User {
// 	return &models.User{
// 		Name:        dto.Name,
// 		FollowUsers: dto.FollowUsers,
// 	}
// }
