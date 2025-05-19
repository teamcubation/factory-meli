package tweet

import (
	"Yuri/twitter/core/ports/services"
	"Yuri/twitter/internal/tweet/handler/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service services.ITweetService
}

func (h *Handler) CreateTweet(c *gin.Context) {
	var req dtos.CreateTweetDTORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}
	err := h.Service.CreateTweet(dtos.ToModel(&req))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario criado com sucesso"})
}
