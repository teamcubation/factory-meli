package handlers

import (
	"Yuri/twitter/core/ports/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	Service services.ITimelineService
}

func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do usuário não pode ser vazio"})
		return
	}
	timeline, err := h.Service.GetTimeline(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeline)
}
