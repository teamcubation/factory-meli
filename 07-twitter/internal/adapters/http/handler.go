package http

import (
	service "07-twitter/core/ports/services"
	"07-twitter/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *Handler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id é obrigatório"})
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	user, err := h.UserService.GetUserById(id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
