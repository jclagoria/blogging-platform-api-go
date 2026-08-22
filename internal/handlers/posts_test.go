package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/juanka/blogging-platform-api/internal/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(handler *PostsHandler) *chi.Mux {
	r := chi.NewRouter()
	generated.HandlerFromMux(handler, r)
	return r
}

func TestCreatePost_Success(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	body := `{"title":"Test Post","content":"Test content","category":"Tech","tags":["go"]}`
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var post generated.Post
	err := json.NewDecoder(w.Body).Decode(&post)
	require.NoError(t, err)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "Test content", post.Content)
	assert.Equal(t, "Tech", post.Category)
	assert.Equal(t, []string{"go"}, *post.Tags)
	assert.NotZero(t, post.CreatedAt)
	assert.Equal(t, post.CreatedAt, post.UpdatedAt)
}

func TestCreatePost_MissingTitle(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	body := `{"content":"Test content","category":"Tech"}`
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var problem ProblemDetail
	err := json.NewDecoder(w.Body).Decode(&problem)
	require.NoError(t, err)
	assert.Equal(t, "/problems/validation-error", problem.Type)
	assert.Len(t, problem.Errors, 1)
	assert.Equal(t, "title", problem.Errors[0].Field)
	assert.Equal(t, "REQUIRED", problem.Errors[0].Code)
}

func TestCreatePost_MissingContent(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	body := `{"title":"Test","category":"Tech"}`
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreatePost_MissingCategory(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	body := `{"title":"Test","content":"content"}`
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestGetPost_Success(t *testing.T) {
	store := NewInMemoryPostStore()
	_, err := store.CreatePost("Test", "Content", "Tech", nil)
	require.NoError(t, err)
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var post generated.Post
	err = json.NewDecoder(w.Body).Decode(&post)
	require.NoError(t, err)
	assert.Equal(t, 1, post.Id)
	assert.Equal(t, "Test", post.Title)
}

func TestGetPost_NotFound(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/posts/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var problem ProblemDetail
	err := json.NewDecoder(w.Body).Decode(&problem)
	require.NoError(t, err)
	assert.Equal(t, "/problems/not-found", problem.Type)
}

func TestListPosts_Empty(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []generated.Post
	err := json.NewDecoder(w.Body).Decode(&posts)
	require.NoError(t, err)
	assert.Empty(t, posts)
}

func TestListPosts_WithData(t *testing.T) {
	store := NewInMemoryPostStore()
	_, err := store.CreatePost("Post 1", "Content 1", "Tech", nil)
	require.NoError(t, err)
	_, err = store.CreatePost("Post 2", "Content 2", "Life", nil)
	require.NoError(t, err)
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []generated.Post
	err = json.NewDecoder(w.Body).Decode(&posts)
	require.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestListPosts_FilterByTerm(t *testing.T) {
	store := NewInMemoryPostStore()
	_, err := store.CreatePost("Go Patterns", "Learn Go", "Tech", nil)
	require.NoError(t, err)
	_, err = store.CreatePost("Rust Basics", "Learn Rust", "Tech", nil)
	require.NoError(t, err)
	_, err = store.CreatePost("Go Concurrency", "Advanced Go", "Tech", nil)
	require.NoError(t, err)
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/posts?term=Go", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var posts []generated.Post
	err = json.NewDecoder(w.Body).Decode(&posts)
	require.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestUpdatePost_Success(t *testing.T) {
	store := NewInMemoryPostStore()
	_, err := store.CreatePost("Original", "Content", "Tech", nil)
	require.NoError(t, err)
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	body := `{"title":"Updated","content":"New content","category":"Life"}`
	req := httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var post generated.Post
	err = json.NewDecoder(w.Body).Decode(&post)
	require.NoError(t, err)
	assert.Equal(t, "Updated", post.Title)
	assert.Equal(t, "New content", post.Content)
}

func TestUpdatePost_NotFound(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	body := `{"title":"Updated","content":"content","category":"Tech"}`
	req := httptest.NewRequest(http.MethodPut, "/posts/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePost_Success(t *testing.T) {
	store := NewInMemoryPostStore()
	_, err := store.CreatePost("To Delete", "Content", "Tech", nil)
	require.NoError(t, err)
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeletePost_NotFound(t *testing.T) {
	store := NewInMemoryPostStore()
	handler := &PostsHandler{Store: store}
	r := setupRouter(handler)

	req := httptest.NewRequest(http.MethodDelete, "/posts/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
