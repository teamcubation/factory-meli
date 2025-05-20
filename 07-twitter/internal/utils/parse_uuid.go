package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ParseStringToUuid(value string) (uuid.UUID, error) {
	idUuid, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, errors.New("[uuid.Parse()] - " + err.Error())
	}

	return idUuid, nil
}
