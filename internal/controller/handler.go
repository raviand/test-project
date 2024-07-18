package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/raviand/test-project/internal/audit"
	"github.com/raviand/test-project/internal/service"
)

type Handler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateCreate(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	TokenMiddleware(next http.Handler) http.Handler
	AuditLog(next http.Handler) http.Handler
}

type handler struct {
	service      service.ProductService
	auditChannel chan<- audit.AuditLog
}

func NewHandler(service service.ProductService, auditChannel chan<- audit.AuditLog) Handler {
	return &handler{
		service:      service,
		auditChannel: auditChannel,
	}
}

// TokenMiddleware checks if the request has the correct token
func (h *handler) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		secretToken := os.Getenv("TOKEN")
		if token != secretToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AuditLog logs the request method, path, timestamp, and user
func (h *handler) AuditLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		audit := audit.AuditLog{
			Method:    r.Method,
			Path:      r.URL.Path,
			TimeStamp: time.Now().Format(time.RFC3339),
			User:      r.Header.Get("X-User"),
		}
		h.auditChannel <- audit
		next.ServeHTTP(w, r)
	})
}
