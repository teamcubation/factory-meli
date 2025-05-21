package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseStringToUuid(value string) (uuid.UUID, error) {
	idUuid, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, errors.New("[uuid.Parse()] - " + err.Error())
	}

	return idUuid, nil
}

func ParseUUIDOrAbort(c *gin.Context, value, field string) (uuid.UUID, bool) {
	id, err := uuid.Parse(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": field + " inválido"})
		return uuid.UUID{}, false
	}
	return id, true
}
