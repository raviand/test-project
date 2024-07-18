package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	rt.Route("/product", func(r chi.Router) {
		r.Get("/", handler.GetAll)
		r.Get("/{id}", handler.GetByID)
		r.Post("/", handler.Create)
		r.Put("/{id}", handler.UpdateCreate)
		r.Patch("/{id}", handler.Patch)
		r.Delete("/{id}", handler.Delete)
	})

	return &server{
		rt: rt,
	}
}

func (s *server) Start() error {
	return http.ListenAndServe(":8080", s.rt)
}
