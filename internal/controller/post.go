package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/raviand/test-project/pkg"
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var product pkg.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	t, err := time.Parse("02/01/2006", product.Expiration)
	if err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.WrongFieldValues))
		return
	}
	newProduct := &pkg.Product{
		Name:        product.Name,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CodeValue:   product.CodeValue,
		IsPublished: product.IsPublished,
		Expiration:  t,
	}
	if err := h.service.Create(newProduct); err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}

	pkg.BuildOkResponse(w, newProduct, http.StatusCreated)
}
