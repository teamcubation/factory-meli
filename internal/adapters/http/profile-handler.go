package http

import (
	"net/http"
	"tweet-tq-rafael/core/dtos"
	"tweet-tq-rafael/core/ports"
	"tweet-tq-rafael/internal/adapters/http/presenters"
	"tweet-tq-rafael/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	Service ports.UserService
}

func (h *UserHandler) Create(c *gin.Context) {
	var createUserRequestDTO dtos.CreateUserRequestDTO
	if err := c.BindJSON(&createUserRequestDTO); err != nil {
		c.JSON(http.StatusBadRequest, "Error on body content.")
		return
	}

	user, err := h.Service.Create(createUserRequestDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error.")
	}

	var response presenters.UserPresenter
	utils.MapSharedFields(user, &response)

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) Follow(c *gin.Context) {
	var followRequestDTO dtos.FollowRequestDTO
	if err := c.BindJSON(&followRequestDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	if err := h.Service.Follow(followRequestDTO); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "follow created"})
}

func (h *UserHandler) Unfollow(c *gin.Context) {

	var unfollowRequestDTO dtos.UnfollowRequestDTO
	if err := c.BindJSON(&unfollowRequestDTO); err != nil {
		c.JSON(http.StatusBadRequest, "invalid input")
		return
	}

	err := h.Service.Unfollow(unfollowRequestDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollow realized"})
}

func (h *UserHandler) Timeline(c *gin.Context) {

	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid input")
	}

	timelineRequestDTO := dtos.TimelineRequestDTO{UserID: id}

	tweets, err := h.Service.Timeline(timelineRequestDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []presenters.TweetPresenter

	for i := 0; i < len(tweets); i++ {
		var item presenters.TweetPresenter
		utils.MapSharedFields(tweets[i], &item)
		response = append(response, item)
	}

	c.JSON(http.StatusOK, response)
}
