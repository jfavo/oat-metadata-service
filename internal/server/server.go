package server

import (
	"log/slog"
	"os"

	"github.com/jfavo/oat-metadata-service/internal/config"
	"github.com/jfavo/oat-metadata-service/internal/services"
)

func StartServer() error {
	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Grab environment variables that we need
	env := config.InitializeEnvironmentWithDefaults(map[string]string{
		"HTTP_PORT": "8080",
	})

	service := services.CreateService(logger)

	service.StartServer(env.GetVariable("HTTP_PORT"))

	return nil
}