package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s service) GetProductByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)

	}
}

func (s service) GetProductById(w http.ResponseWriter, r *http.Request) {

}
