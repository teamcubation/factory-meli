package handlers

import (
	service "07-twitter/core/ports/services"
	"07-twitter/internal/services"
	"07-twitter/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler agrupa os serviços usados pelos endpoints HTTP.
type UserHandler struct {
	UserService service.UserService
}

// NewHandler() cria uma nova instância de Handler.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

type createUserRequest struct {
	Username string `json:"username"`
}

// CreateUserHandler() lida com a criação de um novo usuário.
// @route POST /users
func (h *UserHandler) CreateUserHandler(c *gin.Context) {
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

func (h *UserHandler) GetUserById(c *gin.Context) {
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

// TODO - Quero saber onde posso armazenar isso ao invés de ficar dentro do próprio handler
type followRequest struct {
	UserID      string `json:"user_id"`
	FollowingID string `json:"following_id"`
}

func (h *UserHandler) Following(c *gin.Context) {
	var req followRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if req.UserID == "" || req.FollowingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id e following_id são obrigatórios"})
		return
	}

	userIdUuid, err := utils.ParseStringToUuid(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id inválido"})
		return
	}
	followingIdUuuid, err := utils.ParseStringToUuid(req.FollowingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "following_id inválido"})
		return
	}

	if err := h.UserService.Following(userIdUuid, followingIdUuuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seguindo com sucesso"})
}
