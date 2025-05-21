package dtos

import (
	"Yuri/twitter/core/models"
	"time"
)

type TimelineResponseDTO struct {
	Timeline []*tweetTimeline `json:"timeline"`
}

type tweetTimeline struct {
	ID          string    `json:"_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required,max=280"`
	AuthorId    string    `json:"author_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

func ModelToDTO(model []*models.Tweet) *TimelineResponseDTO {
	var tweetsTimeline []*tweetTimeline
	for _, tweet := range model {
		tweetsTimeline = append(tweetsTimeline, &tweetTimeline{
			ID:          tweet.ID,
			Title:       tweet.Title,
			Description: tweet.Description,
			AuthorId:    tweet.AuthorId,
			CreatedAt:   tweet.CreatedAt,
		})
	}
	return &TimelineResponseDTO{
		Timeline: tweetsTimeline,
	}
}
