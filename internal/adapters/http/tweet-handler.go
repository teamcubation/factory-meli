package http

import (
	"net/http"
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/models"
	"tweet-tq-rafael/core/ports"
	"tweet-tq-rafael/internal/adapters/http/presenters"
	"tweet-tq-rafael/utils"

	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	Service ports.TweetService
}

func (h *TweetHandler) Create(c *gin.Context) {
	var createTweetRequestDTO dtos.CreateTweetRequestDTO
	if err := c.BindJSON(&createTweetRequestDTO); err != nil || len(createTweetRequestDTO.Content) > models.MAX_TWEET_LEN {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tweet, err := h.Service.Create(createTweetRequestDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var response presenters.TweetPresenter
	utils.MapSharedFields(tweet, &response)
	c.JSON(http.StatusOK, response)
}
