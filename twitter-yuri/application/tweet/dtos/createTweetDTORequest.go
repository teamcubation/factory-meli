package dtos

import (
	"Yuri/twitter/core/models"
	"time"
)

type CreateTweetDTORequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required,max=280"`
	AuthorId    string `json:"author_id" binding:"required"`
}

func ToModel(dto *CreateTweetDTORequest) *models.Tweet {
	return &models.Tweet{
		Title:       dto.Title,
		Description: dto.Description,
		AuthorId:    dto.AuthorId,
		CreatedAt:   time.Now(),
	}
}
