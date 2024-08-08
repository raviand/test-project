package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raviand/test-project/internal/controller"
)

type Interface interface {
	Start() error
}

type server struct {
	rt *chi.Mux
}

func NewServer(handler controller.Handler) Interface {
	rt := chi.NewRouter()
	// token middleware to check if the request has the correct token
	rt.Use(middleware.DefaultLogger)
	rt.Use(handler.TokenMiddleware)
	rt.Use(handler.AuditLog)

	rt.Route("/product", func(r chi.Router) {
		r.Get("/", handler.GetAll)
		r.Get("/{id}", handler.GetByID)
		r.Post("/", handler.Create)
		r.Put("/{id}", handler.UpdateCreate)
		r.Patch("/{id}", handler.Patch)
		r.Delete("/{id}", handler.Delete)
	})

	rt.Route("/user", func(r chi.Router) {
		r.Get("/{id}", handler.GetUserById)
		r.Post("/", handler.CreateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})

	return &server{
		rt: rt,
	}
}

func (s *server) Start() error {
	return http.ListenAndServe(":8080", s.rt)
}
