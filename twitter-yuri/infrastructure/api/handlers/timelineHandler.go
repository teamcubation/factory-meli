package handlers

import (
	dtos "Yuri/twitter/application/timeline/dto"
	"Yuri/twitter/core/ports/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	Service services.ITimelineService
}

func (h *TimelineHandler) GetTimeline(c *gin.Context) {
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

	timeline, err := h.Service.GetTimeline(userId, page, tweetsPerPage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.ModelToDTO(timeline))
}
