package utils

import (
	"07-twitter/internal/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindAndValidateFollowRequest(c *gin.Context) (dtos.FollowRequest, bool) {
	var req dtos.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return req, false
	}

	if req.UserID == "" || req.FollowingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id e following_id são obrigatórios"})
		return req, false
	}
	
	return req, true
}
