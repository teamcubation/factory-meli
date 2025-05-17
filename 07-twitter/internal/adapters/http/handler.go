package http

import (
	service "07-twitter/core/ports/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler agrupa os serviços usados pelos endpoints HTTP.
type Handler struct {
	UserService service.UserService
}

// NewHandler() cria uma nova instância de Handler.
func NewHandler(userService service.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

type createUserRequest struct {
	Username string `json:"username"`
}

// CreateUserHandler() lida com a criação de um novo usuário.
// @route POST /users
func (h *Handler) CreateUserHandler(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	newUser, err := h.UserService.CreateUser(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
