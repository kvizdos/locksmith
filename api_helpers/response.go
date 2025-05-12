package api_helpers

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter) {
	WriteResponse(w, APIGenericSuccess, http.StatusOK)
}

func WriteResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	js, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
}

type APIResponseError struct {
	// Human-defined error message
	Reason string `json:"reason,omitempty"`
	// Optional error directly from Go output // extra details
	Details string `json:"details,omitempty"`
}

var APIGenericSuccess = map[string]interface{}{
	"success": true,
}

var APIGenericBadMethod = APIResponseError{
	Reason: "method not allowed",
}
