package user

import (
	"Yuri/twitter/core/ports/services"
	"Yuri/twitter/internal/user/handler/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service services.IUserService
}

func (h *Handler) CreateUser(c *gin.Context) {
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

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario criado com sucesso"})
}
