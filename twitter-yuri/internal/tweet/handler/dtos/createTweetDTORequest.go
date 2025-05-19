package dtos

import (
	"Yuri/twitter/core/models"
)

type CreateTweetDTORequest struct {
	Title       string      `json:"title" binding:"required"`
	Description string      `json:"description" binding:"required"`
	Author      models.User `json:"author" binding:"required"`
}

func ToModel(dto *CreateTweetDTORequest) *models.Tweet {
	return &models.Tweet{
		Title:       dto.Title,
		Description: dto.Description,
		Author:      dto.Author,
	}
}
