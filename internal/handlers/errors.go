package handlers

import (
	"encoding/json"
	"net/http"
)

// ProblemDetail represents an RFC 7807 Problem Details error.
type ProblemDetail struct {
	Type   string               `json:"type"`
	Title  string               `json:"title"`
	Status int                  `json:"status"`
	Detail string               `json:"detail"`
	Errors []ValidationError    `json:"errors,omitempty"`
}

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func writeError(w http.ResponseWriter, status int, problemType, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ProblemDetail{
		Type:   problemType,
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

func writeValidationError(w http.ResponseWriter, errors []ValidationError) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(ProblemDetail{
		Type:   "/problems/validation-error",
		Title:  "Validation Error",
		Status: http.StatusUnprocessableEntity,
		Detail: "Request body failed validation",
		Errors: errors,
	})
}

func writeNotFound(w http.ResponseWriter, detail string) {
	writeError(w, http.StatusNotFound, "/problems/not-found", "Not Found", detail)
}

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "/problems/internal-error", "Internal Server Error", "An unexpected error occurred")
}
