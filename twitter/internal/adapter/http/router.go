package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/twitter-tq/vinofsteel/internal/services"
)

// Router represents the HTTP router for the application
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router with the given services
type NewRouterParams struct {
	UserService services.UserServiceImpl
}

func NewRouter(params NewRouterParams) *Router {
	mux := http.NewServeMux()

	// User routes
	userHandlers := UserHandlers{params.UserService}
	mux.HandleFunc("POST /users", userHandlers.CreateUser)

	return &Router{
		mux: mux,
	}
}

// Run starts the HTTP server on the given address
func (r *Router) Run(ctx context.Context, address string) error {
	slog.InfoContext(ctx, "Starting HTTP server", "address", address)
	return http.ListenAndServe(address, r.mux)
}
