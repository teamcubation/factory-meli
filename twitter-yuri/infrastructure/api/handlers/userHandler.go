package handlers

import (
	"Yuri/twitter/application/user/dtos"
	"Yuri/twitter/core/ports/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service services.IUserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dtos.CreateUserRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	err := h.Service.CreateUser(dtos.ToModel(&req))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully to created user"})
}

func (h *UserHandler) FollowUser(c *gin.Context) {
	var req dtos.UserActionRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	err := h.Service.FollowUser(req.UserID, req.TargetUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully to follow the user"})
}

func (h *UserHandler) UnfollowUser(c *gin.Context) {
	var req dtos.UserActionRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	err := h.Service.UnfollowUser(req.UserID, req.TargetUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully to unfollow the user"})
}
