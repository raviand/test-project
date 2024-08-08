package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/pkg"
)

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}
	pkg.BuildOkResponse(w, products, http.StatusOK)
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	product, err := h.service.GetByID(id)
	if err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}
	pkg.BuildOkResponse(w, product, http.StatusOK)
}

func (h *handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}
	if user == nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.NotFound))
		return
	}
	pkg.BuildOkResponse(w, user, http.StatusOK)
}
