package http

import (
	"net/http"
	"strconv"
	"twitter-api-clone/core/models"
	ports "twitter-api-clone/core/ports"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	TweetService  ports.TweetService
	UserService   ports.UserService
	FollowService ports.FollowService
}

func NewRouter(TweetService ports.TweetService, UserService ports.UserService, FollowService ports.FollowService) *gin.Engine {
	h := &Handler{TweetService: TweetService, UserService: UserService, FollowService: FollowService}
	router := gin.Default()

	router.GET("/healthcheck", h.HealthCheck)

	router.POST("/tweets", h.CreateTweet)
	router.GET("/users/:id/tweets", h.ListTweets)

	router.POST("/users", h.CreateUser)
	router.DELETE("/users/:id", h.DeleteUser)

	router.POST("/follow", h.FollowUser)
	router.POST("/unfollow", h.UnfollowUser)

	router.GET("/timeline", h.GetTimeline)

	return router
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

// TODO: ajustar o DTO pois o nome dos parametros atualmente estao muito especificos para uso interno do projeto (talvez associar nomes json)
//
//	INPUT
//	{
//	  "user_id": "123",
//	  "content": "This is my first tweet from curl 🚀"
//	}
//
// OUTPUT
//
//	{
//	    "Id": "34c0f820-6519-48b7-8175-9263146c9722",
//	    "UserId": "7e30ed3f-89e8-4d7b-8f53-0c163a114163",
//	    "Content": "Tweet Pedro 5",
//	    "CreatedAt": "2025-05-27T13:17:22.5854136-03:00"
//	}
func (h *Handler) CreateTweet(c *gin.Context) {
	var tweet models.Tweet

	if err := c.ShouldBindJSON(&tweet); err != nil {
		return
	}

	if err := h.TweetService.Create(&tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tweet)
}

// TODO: pensar no retorno quando n tiver nenhum tweet para o usuario X, atualmente retornando null
// TODO: vamos usar incremental id ou uuid para o userId?
func (h *Handler) ListTweets(c *gin.Context) {
	userId := c.Param("id")

	tweets, err := h.TweetService.ListByUser(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweets)
}

// TODO: ajustar retorno
// TODO: ajustar para evitar usuario duplicado
// TODO: adicionar outros campos unicos, como email
//
//	{
//	    "Id": "04bbbc90-f8d6-413e-b024-8a1c3ff0c219",
//	    "Name": "Theodoro",
//	    "CreatedAt": "2025-05-27T13:15:59.8488358-03:00"
//	}
func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		return
	}

	if err := h.UserService.Create(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userId := c.Param("id")

	if err := h.UserService.Delete(userId); err != nil {
		return
	}

	c.JSON(http.StatusOK, "user deleted")
}

func (h *Handler) FollowUser(c *gin.Context) {
	var follow models.Follow

	if err := c.ShouldBindJSON(&follow); err != nil {
		return
	}

	if err := h.FollowService.Follow(follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "")
}

func (h *Handler) UnfollowUser(c *gin.Context) {
	follower := c.Param("followerId")
	followed := c.Param("followedId")

	follow := models.Follow{
		FollowedId: followed,
		FollowerId: follower,
	}

	if err := h.FollowService.Unfollow(follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "")
}

func (h *Handler) GetTimeline(c *gin.Context) {
	userId := c.Query("user_id")
	limit, err := strconv.Atoi(c.Query("limit"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweets, err := h.TweetService.GetTimeline(userId, limit, page)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tweets)
}
