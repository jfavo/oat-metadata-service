package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/render"
)

type IMetadataService interface {
	CreateService() IMetadataService
	StartServer(port string)

	health(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
}

type MetadataService struct {
	IMetadataService

	httpRouter 	chi.Router
	logger 		*slog.Logger
}

func CreateService(logger *slog.Logger) IMetadataService {

	router := chi.NewRouter()

	service := MetadataService{}
	service.logger = logger

	// Set middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// Setup route handlers
	router.Get("/health", service.health)
	router.Get("/users/{userId}", service.GetUserById)

	logger.Info("Available routes:")
	for _, route := range router.Routes() {
		logger.Info(route.Pattern)
	}

	service.httpRouter = router

	return service
}

func (ms MetadataService) StartServer(port string) {
	server := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		Handler: ms.httpRouter,
	};

	// Server run context object
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrup/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30 * time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				ms.logger.Error("graceful shutdown timed out... forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			ms.logger.Error(err.Error())
		}
		ms.logger.Info("Server shutting down.")
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		ms.logger.Error(err.Error())
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func (ms MetadataService) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
