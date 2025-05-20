package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Setups slog for the application
func SetupLogger(ctx context.Context) {
	env := os.Getenv("ENV")

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	// #nosec G301
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		slog.ErrorContext(ctx, "Error creating logs directory", "error", err)
		os.Exit(1)
	}

	// Create log file with timestamp in filename
	timestamp := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logsDir, fmt.Sprintf("%s.log", timestamp))

	// #nosec G304 G302
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.ErrorContext(ctx, "Error opening log file", "error", err)
		os.Exit(1)
	}

	// Create a multi-writer that writes to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var handler slog.Handler

	switch env {
	case "production":
		// Use JSON format in production for better parsing by log aggregation tools
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			// Add timestamp with consistent ISO8601 format
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
				}
				return a
			},
		})
	case "test":
		// Minimal logging in test environment
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})
	default:
		// Development environment with more verbose logging
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true, // Include source file and line in logs for development
		})
	}

	// Replace the default slog logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Log the initial configuration
	slog.InfoContext(ctx, "Logger initialized",
		"environment", env,
		"level", logger.Handler().Enabled(ctx, slog.LevelInfo),
		"logFile", logFilePath,
	)
}
