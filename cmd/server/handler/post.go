package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raviand/test-project/internal/domain"
)

func (s service) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("My-Header"))
	var product domain.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		setResponse(w, http.StatusBadRequest, nil, false, err)
		return
	}
	productoDeBaseDeDatos, err := s.db.CreateProduct(&product)
	if err != nil {
		setResponse(w, http.StatusInternalServerError, nil, false, err)
		return
	}
	fmt.Println(productoDeBaseDeDatos)
	setResponse(w, http.StatusCreated, productoDeBaseDeDatos, true, nil)
}
