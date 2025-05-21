package handlers

import (
	"Yuri/twitter/application/tweet/dtos"
	"Yuri/twitter/core/ports/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	Service services.ITweetService
}

func (h *TweetHandler) CreateTweet(c *gin.Context) {
	var req dtos.CreateTweetDTORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Service.CreateTweet(dtos.ToModel(&req))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully to created tweet"})
}

func (h *TweetHandler) FindTweetsByUserId(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty user ID"})
		return
	}
	tweets, err := h.Service.ListTweetsByUserID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, dtos.ModelToDTO(tweets))
}
