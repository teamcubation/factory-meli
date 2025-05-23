package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	coremodels "07-twitter/core/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct {
	CreateUserFunc  func(string) (*coremodels.User, error)
	GetUserByIdFunc func(string) (*coremodels.User, error)
	FollowFunc      func(uuid.UUID, uuid.UUID) error
	UnfollowFunc    func(uuid.UUID, uuid.UUID) error
}

func (m *mockUserService) CreateUser(username string) (*coremodels.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(username)
	}
	return nil, nil
}
func (m *mockUserService) GetUserById(id string) (*coremodels.User, error) {
	if m.GetUserByIdFunc != nil {
		return m.GetUserByIdFunc(id)
	}
	return nil, nil
}
func (m *mockUserService) Follow(userId, followingId uuid.UUID) error {
	if m.FollowFunc != nil {
		return m.FollowFunc(userId, followingId)
	}
	return nil
}
func (m *mockUserService) Unfollow(userId, followingId uuid.UUID) error {
	if m.UnfollowFunc != nil {
		return m.UnfollowFunc(userId, followingId)
	}
	return nil
}

func TestCreateUserHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"username":"testuser"}`)
	c.Request = httptest.NewRequest("POST", "/users", body)
	c.Request.Header.Set("Content-Type", "application/json")

	h := &UserHandler{UserService: &mockUserService{
		CreateUserFunc: func(username string) (*coremodels.User, error) {
			return &coremodels.User{ID: uuid.New(), Username: username}, nil
		},
	}}

	h.CreateUserHandler(c)
	assert.Equal(t, 201, w.Code)
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString("{invalid-json}")
	c.Request = httptest.NewRequest("POST", "/users", body)
	c.Request.Header.Set("Content-Type", "application/json")

	h := &UserHandler{UserService: &mockUserService{}}
	h.CreateUserHandler(c)
	assert.Equal(t, 400, w.Code)
}

func TestGetUserById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id := uuid.New().String()
	c.Params = gin.Params{{Key: "id", Value: id}}

	h := &UserHandler{UserService: &mockUserService{
		GetUserByIdFunc: func(id string) (*coremodels.User, error) {
			return &coremodels.User{ID: uuid.MustParse(id), Username: "testuser"}, nil
		},
	}}

	h.GetUserById(c)
	assert.Equal(t, 200, w.Code)
}

func TestGetUserById_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	h := &UserHandler{UserService: &mockUserService{}}
	h.GetUserById(c)
	assert.Equal(t, 400, w.Code)
}

func TestFollow_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"user_id":"` + uuid.New().String() + `", "following_id":"` + uuid.New().String() + `"}`)
	c.Request = httptest.NewRequest("POST", "/users/follow", body)
	c.Request.Header.Set("Content-Type", "application/json")

	h := &UserHandler{UserService: &mockUserService{
		FollowFunc: func(userId, followingId uuid.UUID) error { return nil },
	}}

	h.Follow(c)
	assert.Equal(t, 200, w.Code)
}

func TestUnfollow_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"user_id":"` + uuid.New().String() + `", "following_id":"` + uuid.New().String() + `"}`)
	c.Request = httptest.NewRequest("POST", "/users/unfollow", body)
	c.Request.Header.Set("Content-Type", "application/json")

	h := &UserHandler{UserService: &mockUserService{
		UnfollowFunc: func(userId, followingId uuid.UUID) error { return nil },
	}}

	h.Unfollow(c)
	assert.Equal(t, 200, w.Code)
}
