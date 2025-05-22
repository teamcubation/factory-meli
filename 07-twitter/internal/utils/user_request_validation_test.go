package utils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"07-twitter/internal/dtos"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_BindAndValidateFollowRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := dtos.FollowRequest{
		UserID:      "user-uuid",
		FollowingID: "follow-uuid",
	}
	jsonBody, _ := json.Marshal(body)
	c.Request = httptest.NewRequest("POST", "/follow", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	result, ok := BindAndValidateFollowRequest(c)

	assert.True(t, ok)
	assert.Equal(t, body, result)
}

func Test_BindAndValidateFollowRequest_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", "/follow", bytes.NewBuffer([]byte("{invalid-json}")))
	c.Request.Header.Set("Content-Type", "application/json")

	result, ok := BindAndValidateFollowRequest(c)

	assert.False(t, ok)
	assert.Equal(t, dtos.FollowRequest{}, result)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"JSON inválido"}`, w.Body.String())
}

func Test_BindAndValidateFollowRequest_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := dtos.FollowRequest{UserID: "", FollowingID: ""}
	jsonBody, _ := json.Marshal(body)
	c.Request = httptest.NewRequest("POST", "/follow", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	result, ok := BindAndValidateFollowRequest(c)

	assert.False(t, ok)
	assert.Equal(t, dtos.FollowRequest{}, result)
	assert.Equal(t, 400, w.Code)
	assert.JSONEq(t, `{"error":"user_id e following_id são obrigatórios"}`, w.Body.String())
}
