package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/juanka/blogging-platform-api/internal/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 24-char zero hex string — valid ObjectID format but no matching document
var nonExistentID = fmt.Sprintf("%024x", 0)

func TestErrorFormat_NotFound(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := chi.NewRouter()
	generated.HandlerFromMux(handler, r)

	req := httptest.NewRequest(http.MethodGet, "/posts/"+nonExistentID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

	var problem ProblemDetail
	err := json.NewDecoder(w.Body).Decode(&problem)
	require.NoError(t, err)
	assert.Equal(t, "/problems/not-found", problem.Type)
	assert.Equal(t, "Not Found", problem.Title)
	assert.Equal(t, 404, problem.Status)
	assert.NotEmpty(t, problem.Detail)
	assert.Nil(t, problem.Errors)
}

func TestErrorFormat_ValidationErrors(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := chi.NewRouter()
	generated.HandlerFromMux(handler, r)

	body := `{"title":"","content":"","category":""}`
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var problem ProblemDetail
	err := json.NewDecoder(w.Body).Decode(&problem)
	require.NoError(t, err)
	assert.Equal(t, "/problems/validation-error", problem.Type)
	assert.Equal(t, "Validation Error", problem.Title)
	assert.Equal(t, 422, problem.Status)
	assert.NotEmpty(t, problem.Detail)
	require.Len(t, problem.Errors, 3)

	fields := make(map[string]bool)
	for _, e := range problem.Errors {
		assert.Equal(t, "REQUIRED", e.Code)
		assert.NotEmpty(t, e.Message)
		fields[e.Field] = true
	}
	assert.True(t, fields["title"])
	assert.True(t, fields["content"])
	assert.True(t, fields["category"])
}

func TestErrorFormat_UpdateNotFound(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := chi.NewRouter()
	generated.HandlerFromMux(handler, r)

	body := `{"title":"T","content":"C","category":"Cat"}`
	req := httptest.NewRequest(http.MethodPut, "/posts/"+nonExistentID, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotFound, w.Code)

	var problem ProblemDetail
	err := json.NewDecoder(w.Body).Decode(&problem)
	require.NoError(t, err)
	assert.Equal(t, "/problems/not-found", problem.Type)
}

func TestErrorFormat_DeleteNotFound(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := chi.NewRouter()
	generated.HandlerFromMux(handler, r)

	req := httptest.NewRequest(http.MethodDelete, "/posts/"+nonExistentID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotFound, w.Code)

	var problem ProblemDetail
	err := json.NewDecoder(w.Body).Decode(&problem)
	require.NoError(t, err)
	assert.Equal(t, "/problems/not-found", problem.Type)
}
