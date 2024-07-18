package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/pkg"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		pkg.BuildErrorResponse(w, pkg.GetError(pkg.BadRequest))
		return
	}
	if err := h.service.Delete(id); err != nil {
		pkg.BuildErrorResponse(w, err.(pkg.ApiError))
		return
	}
	pkg.BuildOkResponse(w, nil, http.StatusNoContent)
}
