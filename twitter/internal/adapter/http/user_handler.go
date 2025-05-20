package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/internal/services"
)

// UserHandler handles HTTP requests for user operations
type UserHandlers struct {
	service services.UserServiceImpl
}

// Internal user type to return json correctly
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	ApiKey    string    `json:"api_key"`
}

func (h UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type parameters struct {
		Email string `json:"email"`
	}

	// Parse request body
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(ctx, w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Create user using the service
	user, err := h.service.CreateUser(ctx, params.Email)
	if err != nil {
		respondWithError(ctx, w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	// Return success response
	respondWithJSON(ctx, w, http.StatusCreated, modelUserToUser(*user))
}

// Utilities
func modelUserToUser(user models.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}
