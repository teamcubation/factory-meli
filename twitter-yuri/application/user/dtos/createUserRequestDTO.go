package dtos

import (
	"Yuri/twitter/core/models"
)

type CreateUserRequestDTO struct {
	Name        string   `json:"name" validate:"required"`
	FollowUsers []string `json:"follow_users"`
}

// TODO: nessa camada não pode ter dependência do mongoDB?
func ToModel(dto *CreateUserRequestDTO) *models.User {
	return &models.User{
		Name:        dto.Name,
		FollowUsers: dto.FollowUsers,
	}
}
