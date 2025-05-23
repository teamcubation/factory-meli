package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"07-twitter/core/models"
	"07-twitter/internal/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTweetService struct {
	CreateTweetFunc           func(string, uuid.UUID) (*models.Tweet, error)
	GetTweetByUserIdFunc      func(uuid.UUID) ([]*models.Tweet, error)
	GetTimelineByUserIdFunc   func(uuid.UUID, int, int) ([]*models.Tweet, error)
}

func (m *mockTweetService) CreateTweet(content string, userID uuid.UUID) (*models.Tweet, error) {
	return m.CreateTweetFunc(content, userID)
}
func (m *mockTweetService) GetTweetByUserId(userID uuid.UUID) ([]*models.Tweet, error) {
	return m.GetTweetByUserIdFunc(userID)
}
func (m *mockTweetService) GetTimelineByUserId(userID uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
	return m.GetTimelineByUserIdFunc(userID, limit, offset)
}

func newTestGinContextWithJSON(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	c.Request = httptest.NewRequest(method, url, &buf)
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestCreateTweet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	mockService := &mockTweetService{
		CreateTweetFunc: func(content string, uid uuid.UUID) (*models.Tweet, error) {
			return &models.Tweet{ID: uuid.New(), UserID: uid, Content: content}, nil
		},
	}
	h := NewTweetHandler(mockService)
	body := dtos.TweetRequest{UserID: userID.String(), Content: "tweet!"}
	c, w := newTestGinContextWithJSON("POST", "/tweets", body)

	h.CreateTweet(c)
	assert.Equal(t, 201, w.Code)
}

func TestCreateTweet_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockTweetService{}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/tweets", bytes.NewBuffer([]byte("{invalid-json}")))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateTweet(c)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"JSON inválido"}`, w.Body.String())
}

func TestCreateTweet_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockTweetService{}
	h := NewTweetHandler(mockService)
	body := dtos.TweetRequest{UserID: "", Content: ""}
	c, w := newTestGinContextWithJSON("POST", "/tweets", body)

	h.CreateTweet(c)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"Conteúdo e user_id são obrigatórios"}`, w.Body.String())
}

func TestCreateTweet_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockTweetService{}
	h := NewTweetHandler(mockService)
	body := dtos.TweetRequest{UserID: "not-a-uuid", Content: "tweet!"}
	c, w := newTestGinContextWithJSON("POST", "/tweets", body)

	h.CreateTweet(c)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"user_id inválido"}`, w.Body.String())
}

func TestCreateTweet_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	mockService := &mockTweetService{
		CreateTweetFunc: func(content string, uid uuid.UUID) (*models.Tweet, error) {
			return nil, errors.New("erro ao criar tweet")
		},
	}
	h := NewTweetHandler(mockService)
	body := dtos.TweetRequest{UserID: userID.String(), Content: "tweet!"}
	c, w := newTestGinContextWithJSON("POST", "/tweets", body)

	h.CreateTweet(c)
	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error":"erro ao criar tweet"}`, w.Body.String())
}

func TestGetTweetsByUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	mockService := &mockTweetService{
		GetTweetByUserIdFunc: func(uid uuid.UUID) ([]*models.Tweet, error) {
			return []*models.Tweet{{ID: uuid.New(), UserID: uid, Content: "tweet!"}}, nil
		},
	}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	h.GetTweetsByUser(c)
	assert.Equal(t, 200, w.Code)
}

func TestGetTweetsByUser_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockTweetService{}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	h.GetTweetsByUser(c)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"user_id inválido"}`, w.Body.String())
}

func TestGetTweetsByUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	mockService := &mockTweetService{
		GetTweetByUserIdFunc: func(uid uuid.UUID) ([]*models.Tweet, error) {
			return nil, errors.New("erro ao buscar tweets")
		},
	}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	h.GetTweetsByUser(c)
	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error":"erro ao buscar tweets"}`, w.Body.String())
}

func TestGetTweetTimelineByUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	mockService := &mockTweetService{
		GetTimelineByUserIdFunc: func(uid uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
			return []*models.Tweet{{ID: uuid.New(), UserID: uid, Content: "tweet!"}}, nil
		},
	}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	h.GetTweetTimelineByUser(c)
	assert.Equal(t, 200, w.Code)
}

func TestGetTweetTimelineByUser_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &mockTweetService{}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	h.GetTweetTimelineByUser(c)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"user_id inválido"}`, w.Body.String())
}

func TestGetTweetTimelineByUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	mockService := &mockTweetService{
		GetTimelineByUserIdFunc: func(uid uuid.UUID, limit, offset int) ([]*models.Tweet, error) {
			return nil, errors.New("erro ao buscar timeline")
		},
	}
	h := NewTweetHandler(mockService)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: userID.String()}}

	h.GetTweetTimelineByUser(c)
	assert.Equal(t, 500, w.Code)
	assert.JSONEq(t, `{"error":"erro ao buscar timeline"}`, w.Body.String())
}
