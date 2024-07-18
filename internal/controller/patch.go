package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/pkg"
)

func (h *handler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	var product pkg.ProductPatchRequest
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	if err := h.service.Patch(id, &product); err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}
	pkg.BuildOkResponse(w, product, http.StatusOK)
}
