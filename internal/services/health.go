package services

import "net/http"

// Health is an http handler that returns the health of the service
// and any dependent components (DB, auth-service, etc)
func (ms MetadataService) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}