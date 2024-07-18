package pkg

import (
	"encoding/json"
	"net/http"
)

func BuildOkResponse(w http.ResponseWriter, data interface{}, responseCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	json.NewEncoder(w).Encode(data)
}

func BuildErrorResponse(w http.ResponseWriter, err ApiError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.ResponseCode)
	json.NewEncoder(w).Encode(err)
}
