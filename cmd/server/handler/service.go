package handler

import (
	"encoding/json"
	"net/http"

	"github.com/raviand/test-project/internal/repository"
)

type Response struct {
	Success bool        `json:"success"`
	Error   error       `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}

type service struct {
	db repository.DataInterface
}

type Interface interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProductByCode(w http.ResponseWriter, r *http.Request)
	GetProductById(w http.ResponseWriter, r *http.Request)
	PatchProduct(w http.ResponseWriter, r *http.Request)
}

func NewHandler(db repository.DataInterface) Interface {
	return &service{
		db: db,
	}
}

func setResponse(w http.ResponseWriter, responseCode int, data interface{}, success bool, err error) {
	r := Response{
		success,
		err,
		data,
	}
	if b, err := json.Marshal(r); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseCode)
		w.Write(b)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	w.Write(nil)
}
