package services

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
)

// UserContext is an http middleware function
func UserContext(ms *MetadataService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}

func (ms MetadataService) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	ms.logger.Info("GetUserById hit!")

	w.Write([]byte(fmt.Sprintf("GetUserById userId: %s", userId)))
}