package controller

import (
	"net/http"

	"github.com/raviand/test-project/internal/service"
)

type Handler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateCreate(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service service.ProductService
}

func NewHandler(service service.ProductService) Handler {
	return &handler{
		service: service,
	}
}
