package handlers

import (
	service "07-twitter/core/ports/services"
	"07-twitter/internal/dtos"
	"07-twitter/internal/services"
	"07-twitter/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler agrupa os serviços usados pelos endpoints HTTP de usuário.
type UserHandler struct {
	UserService service.UserService
}

// NewUserHandler() cria uma nova instância de UserHandler.
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

// GetUserById() retorna um usuário pelo seu ID.
// @route GET /users/:id
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

// Follow() faz o usuário seguir outro usuário.
// @route POST /users/follow
func (h *UserHandler) Follow(c *gin.Context) {
	var req dtos.FollowRequest

	req, ok := utils.BindAndValidateFollowRequest(c)
	if !ok {
		return
	}

	userId, ok := utils.ParseUUIDOrAbort(c, req.UserID, "user_id")
	if !ok {
		return
	}
	
	followingId, ok := utils.ParseUUIDOrAbort(c, req.FollowingID, "following_id")
	if !ok {
		return
	}

	if err := h.UserService.Follow(userId, followingId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seguindo com sucesso"})
}

// Unfollow() faz o usuário deixar de seguir outro usuário.
// @route POST /users/unfollow
func (h *UserHandler) Unfollow(c *gin.Context) {
	var req dtos.FollowRequest
	req, ok := utils.BindAndValidateFollowRequest(c)
	if !ok {
		return
	}

	userId, ok := utils.ParseUUIDOrAbort(c, req.UserID, "user_id")
	if !ok {
		return
	}

	followingId, ok := utils.ParseUUIDOrAbort(c, req.FollowingID, "following_id")
	if !ok {
		return
	}

	if err := h.UserService.Unfollow(userId, followingId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sucesso ao deixar de seguir usuário"})
}
