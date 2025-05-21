package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/core/domain/models"
	"github.com/twitter-tq/vinofsteel/internal/services"
)

// FollowHandlers handles HTTP requests for follow operations
type FollowHandlers struct {
	service services.FollowServiceImpl
}

// Internal follow type to return json correctly
type Follow struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FollowerID uuid.UUID `json:"follower_id"`
	FollowedID uuid.UUID `json:"followed_id"`
}

func (h FollowHandlers) FollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	followedIDStr := r.PathValue("followed_id")

	type parameters struct {
		FollowerIDStr string `json:"follower_id"`
	}

	// Parse request body
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(ctx, w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Follow user using the service
	follow, err := h.service.FollowUser(ctx, params.FollowerIDStr, followedIDStr)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Return success response
	respondWithJSON(ctx, w, http.StatusCreated, modelFollowToFollow(*follow))
}

// Utilities
func modelFollowToFollow(follow models.Follow) Follow {
	return Follow{
		ID:         follow.ID,
		CreatedAt:  follow.CreatedAt,
		UpdatedAt:  follow.UpdatedAt,
		FollowerID: follow.FollowerID,
		FollowedID: follow.FollowedID,
	}
}

func (h FollowHandlers) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	followedIDStr := r.PathValue("followed_id")

	type parameters struct {
		FollowerIDStr string `json:"follower_id"`
	}

	// Parse request body
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(ctx, w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Unfollow user using the service
	err := h.service.UnfollowUser(ctx, params.FollowerIDStr, followedIDStr)
	if err != nil {
		switch e := err.(type) {
		case services.ServiceError:
			respondWithError(ctx, w, err.(services.ServiceError).Code(), e.Error())
		default:
			respondWithError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Return success response
	respondWithJSON(ctx, w, http.StatusNoContent, "")
}
