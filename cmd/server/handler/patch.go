package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/internal/domain"
)

func (s service) PatchProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		setResponse(w, http.StatusBadRequest, map[string]any{"message": "invalid id", "data": nil}, false, err)
		return
	}
	pr, err := s.db.GetProductById(id)
	if err != nil {
		setResponse(w, http.StatusNotFound, nil, false, err)
		return
	}
	prod := domain.Product{
		Id:          pr.Id,
		Name:        pr.Name,
		Quantity:    pr.Quantity,
		CodeValue:   pr.CodeValue,
		IsPublished: pr.IsPublished,
		Expiration:  pr.Expiration,
		Price:       pr.Price,
	}
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		setResponse(w, http.StatusBadRequest, nil, false, err)
		return
	}
	res, err := s.db.PatchProduct(id, &prod)
	if err != nil {
		setResponse(w, http.StatusInternalServerError, nil, false, err)
		return
	}
	setResponse(w, http.StatusOK, res, true, nil)
}
