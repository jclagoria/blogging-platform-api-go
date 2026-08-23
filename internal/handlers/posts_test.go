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

func setupRouter(handler *PostsHandler) *chi.Mux {
	r := chi.NewRouter()
	generated.HandlerFromMux(handler, r)
	return r
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantTitle  string
		wantField  string
	}{
		{"valid post", `{"title":"Test Post","content":"Test content","category":"Tech","tags":["go"]}`, http.StatusCreated, "Test Post", ""},
		{"missing title", `{"content":"Test content","category":"Tech"}`, http.StatusUnprocessableEntity, "", "title"},
		{"missing content", `{"title":"Test","category":"Tech"}`, http.StatusUnprocessableEntity, "", "content"},
		{"missing category", `{"title":"Test","content":"content"}`, http.StatusUnprocessableEntity, "", "category"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryPostStore()
			handler := &PostsHandler{Store: store}
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusCreated {
				var post generated.Post
				err := json.NewDecoder(w.Body).Decode(&post)
				require.NoError(t, err)
				assert.Equal(t, tt.wantTitle, post.Title)
				assert.NotEmpty(t, post.Id)
				assert.NotZero(t, post.CreatedAt)
				assert.Equal(t, post.CreatedAt, post.UpdatedAt)
			} else {
				var problem ProblemDetail
				err := json.NewDecoder(w.Body).Decode(&problem)
				require.NoError(t, err)
				assert.Equal(t, "/problems/validation-error", problem.Type)
				if tt.wantField != "" {
					assert.Len(t, problem.Errors, 1)
					assert.Equal(t, tt.wantField, problem.Errors[0].Field)
					assert.Equal(t, "REQUIRED", problem.Errors[0].Code)
				}
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	// Create a post first to get a valid ID
	store := NewInMemoryPostStore()
	created, _ := store.CreatePost("Test", "Content", "Tech", nil)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantTitle  string
	}{
		{"existing post", "/posts/" + created.Id, http.StatusOK, "Test"},
		{"non-existent post", "/posts/000000000000000000000000", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PostsHandler{Store: store}
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var post generated.Post
				err := json.NewDecoder(w.Body).Decode(&post)
				require.NoError(t, err)
				assert.Equal(t, tt.wantTitle, post.Title)
				assert.Equal(t, created.Id, post.Id)
			} else {
				var problem ProblemDetail
				err := json.NewDecoder(w.Body).Decode(&problem)
				require.NoError(t, err)
				assert.Equal(t, "/problems/not-found", problem.Type)
			}
		})
	}
}

func TestListPosts(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(store *InMemoryPostStore)
		query      string
		wantStatus int
		wantLen    int
	}{
		{"empty list", func(s *InMemoryPostStore) {}, "", http.StatusOK, 0},
		{"with data", func(s *InMemoryPostStore) {
			_, _ = s.CreatePost("Post 1", "Content 1", "Tech", nil)
			_, _ = s.CreatePost("Post 2", "Content 2", "Life", nil)
		}, "", http.StatusOK, 2},
		{"filter matching", func(s *InMemoryPostStore) {
			_, _ = s.CreatePost("Go Patterns", "Learn Go", "Tech", nil)
			_, _ = s.CreatePost("Rust Basics", "Learn Rust", "Tech", nil)
			_, _ = s.CreatePost("Go Concurrency", "Advanced Go", "Tech", nil)
		}, "?term=Go", http.StatusOK, 2},
		{"filter no match", func(s *InMemoryPostStore) {
			_, _ = s.CreatePost("Go Patterns", "Learn Go", "Tech", nil)
			_, _ = s.CreatePost("Rust Basics", "Learn Rust", "Tech", nil)
		}, "?term=Python", http.StatusOK, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryPostStore()
			tt.setup(store)
			handler := &PostsHandler{Store: store}
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/posts"+tt.query, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var posts []generated.Post
			err := json.NewDecoder(w.Body).Decode(&posts)
			require.NoError(t, err)
			assert.Len(t, posts, tt.wantLen)

			// Verify all returned posts have non-empty IDs
			for _, p := range posts {
				assert.NotEmpty(t, p.Id)
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	store := NewInMemoryPostStore()
	created, _ := store.CreatePost("Original", "Content", "Tech", nil)

	nonExistentID := fmt.Sprintf("%024x", 0) // valid hex but no matching doc

	tests := []struct {
		name       string
		path       string
		body       string
		wantStatus int
		wantTitle  string
	}{
		{"valid update", "/posts/" + created.Id, `{"title":"Updated","content":"New content","category":"Life"}`, http.StatusOK, "Updated"},
		{"non-existent post", "/posts/" + nonExistentID, `{"title":"Updated","content":"content","category":"Tech"}`, http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PostsHandler{Store: store}
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodPut, tt.path, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var post generated.Post
				err := json.NewDecoder(w.Body).Decode(&post)
				require.NoError(t, err)
				assert.Equal(t, tt.wantTitle, post.Title)
				assert.Equal(t, created.Id, post.Id)
			} else {
				var problem ProblemDetail
				err := json.NewDecoder(w.Body).Decode(&problem)
				require.NoError(t, err)
				assert.Equal(t, "/problems/not-found", problem.Type)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	store := NewInMemoryPostStore()
	created, _ := store.CreatePost("To Delete", "Content", "Tech", nil)

	nonExistentID := fmt.Sprintf("%024x", 0)

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"existing post", "/posts/" + created.Id, http.StatusNoContent},
		{"non-existent post", "/posts/" + nonExistentID, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &PostsHandler{Store: store}
			r := setupRouter(handler)

			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusNoContent {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}
