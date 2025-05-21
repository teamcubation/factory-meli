package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// respondWithError sends a JSON error response
func respondWithError(ctx context.Context, w http.ResponseWriter, code int, message string) {
	slog.ErrorContext(ctx, message, "status_code", code)
	respondWithJSON(ctx, w, code, map[string]string{"message": message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	response, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Error marshaling JSON response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
