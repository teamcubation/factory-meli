package handlers

import (
	service "07-twitter/core/ports/services"
	dto "07-twitter/internal/dtos"
	utils "07-twitter/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TweetHandler agrupa os serviços usados pelos endpoints HTTP.
type TweetHandler struct {
	TweetService service.TweetService
}

// NewTweetHandler() cria uma nova instância de Handler.
func NewTweetHandler(tweetService service.TweetService) *TweetHandler {
	return &TweetHandler{
		TweetService: tweetService,
	}
}

// CreateTweetHandler lida com a criação de um novo tweet.
// @route POST /tweets
func (h *TweetHandler) CreateTweet(c *gin.Context) {
	var req dto.TweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if req.Content == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conteúdo e user_id são obrigatórios"})
		return
	}

	userID, err := utils.ParseStringToUuid(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id inválido"})
		return
	}

	tweet, err := h.TweetService.CreateTweet(req.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tweet)
}

// GetTweetsByUserHandler retorna todos os tweets de um usuário específico.
// @route GET /users/:id/tweets
func (h *TweetHandler) GetTweetsByUser(c *gin.Context) {
	id := c.Param("id")
	userID, ok := utils.ParseUUIDOrAbort(c, id, "user_id")
	if !ok {
		return
	}

	tweets, err := h.TweetService.GetTweetByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweets)
}

// GetTimelineHandler retorna os últimos N tweets das contas seguidas pelo usuário.
// @route GET /users/:id/timeline?limit=10&offset=0
func (h *TweetHandler) GetTweetTimelineByUser(c *gin.Context) {
	id := c.Param("id")
	userID, ok := utils.ParseUUIDOrAbort(c, id, "user_id")
	if !ok {
		return
	}

	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	tweets, err := h.TweetService.GetTimelineByUserId(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweets)
}
