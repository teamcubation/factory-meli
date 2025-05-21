package dtos

import (
	"Yuri/twitter/core/models"
	"time"
)

type TweetsByUserDTOResponse struct {
	Data []*object `json:"data"`
}

type object struct {
	ID          string    `json:"_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required,max=280"`
	AuthorId    string    `json:"author_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

func ModelToDTO(model []*models.Tweet) *TweetsByUserDTOResponse {
	var objects []*object
	for _, tweet := range model {
		objects = append(objects, &object{
			ID:          tweet.ID,
			Title:       tweet.Title,
			Description: tweet.Description,
			AuthorId:    tweet.AuthorId,
			CreatedAt:   tweet.CreatedAt,
		})
	}
	return &TweetsByUserDTOResponse{
		Data: objects,
	}
}
