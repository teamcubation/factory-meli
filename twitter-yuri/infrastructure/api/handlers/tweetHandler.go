package handlers

import (
	"Yuri/twitter/application/tweet/dtos"
	"Yuri/twitter/core/ports/services"
	"net/http"
	"strconv"

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
	userId := c.Query("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	page, errPage := strconv.Atoi(pageStr)

	if errPage != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
	}

	tweetsPerPageStr := c.DefaultQuery("tweetsPerPage", "10")
	tweetsPerPage, errLimit := strconv.Atoi(tweetsPerPageStr)

	if errLimit != nil || tweetsPerPage <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
	}

	tweets, err := h.Service.ListTweetsByUserID(userId, page, tweetsPerPage)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.ModelToDTO(tweets))
}
