package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/pkg"
)

func (h *handler) UpdateCreate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	var product pkg.CreateProductRequest
	if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	t, err := time.Parse("02/01/2006", product.Expiration)
	if err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.WrongFieldValues))
		return
	}
	newProduct := &pkg.Product{
		ID:          id,
		Name:        product.Name,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  t,
	}
	if err := h.service.UpdateCreate(newProduct); err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}
	pkg.BuildOkResponse(w, newProduct, http.StatusOK)
}
