package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ParseStringToUuid_Success(t *testing.T) {
	u := uuid.New()

	id, err := ParseStringToUuid(u.String())

	assert.NoError(t, err)
	assert.Equal(t, u, id)
}

func Test_ParseStringToUuid_InvalidUuid(t *testing.T) {
	_, err := ParseStringToUuid("not-a-uuid")

	assert.Error(t, err)
}

func Test_ParseUUIDOrAbort_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	u := uuid.New()

	id, ok := ParseUUIDOrAbort(c, u.String(), "id")

	assert.True(t, ok)
	assert.Equal(t, u, id)
}

func Test_ParseUUIDOrAbort_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	field := "id"

	id, ok := ParseUUIDOrAbort(c, "invalid-uuid", field)

	assert.False(t, ok)
	assert.Equal(t, uuid.UUID{}, id)
	assert.Equal(t, 400, w.Code)
	errorExpected := `{"error":"id inválido"}`
	assert.JSONEq(t, errorExpected, w.Body.String())
}
