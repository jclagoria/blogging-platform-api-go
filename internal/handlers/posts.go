package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/juanka/blogging-platform-api/internal/generated"
)

const (
	contentTypeJSON      = "application/json"
	errPostNotFound      = "The requested post was not found"
)

// PostsHandler implements the ServerInterface for posts endpoints.
type PostsHandler struct {
	Store PostStore
}

// HealthCheck handles GET /health.
func (h *PostsHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	HealthCheck(w, r)
}

// CreatePost handles POST /posts.
func (h *PostsHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req generated.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeInternalError(w)
		return
	}

	var validationErrors []ValidationError
	if strings.TrimSpace(req.Title) == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "title", Message: "title is required", Code: "REQUIRED"})
	}
	if strings.TrimSpace(req.Content) == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "content", Message: "content is required", Code: "REQUIRED"})
	}
	if strings.TrimSpace(req.Category) == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "category", Message: "category is required", Code: "REQUIRED"})
	}
	if len(validationErrors) > 0 {
		writeValidationError(w, validationErrors)
		return
	}

	var tags []string
	if req.Tags != nil {
		tags = *req.Tags
	}

	post, err := h.Store.CreatePost(req.Title, req.Content, req.Category, tags)
	if err != nil {
		writeInternalError(w)
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

// ListPosts handles GET /posts.
func (h *PostsHandler) ListPosts(w http.ResponseWriter, r *http.Request, params generated.ListPostsParams) {
	term := ""
	if params.Term != nil {
		term = *params.Term
	}

	posts, err := h.Store.ListPosts(term)
	if err != nil {
		writeInternalError(w)
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

// GetPost handles GET /posts/{id}.
func (h *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request, id generated.PostId) {
	post, err := h.Store.GetPost(id)
	if err != nil {
		writeNotFound(w, errPostNotFound)
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

// UpdatePost handles PUT /posts/{id}.
func (h *PostsHandler) UpdatePost(w http.ResponseWriter, r *http.Request, id generated.PostId) {
	var req generated.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeInternalError(w)
		return
	}

	var validationErrors []ValidationError
	if strings.TrimSpace(req.Title) == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "title", Message: "title is required", Code: "REQUIRED"})
	}
	if strings.TrimSpace(req.Content) == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "content", Message: "content is required", Code: "REQUIRED"})
	}
	if strings.TrimSpace(req.Category) == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "category", Message: "category is required", Code: "REQUIRED"})
	}
	if len(validationErrors) > 0 {
		writeValidationError(w, validationErrors)
		return
	}

	var tags []string
	if req.Tags != nil {
		tags = *req.Tags
	}

	post, err := h.Store.UpdatePost(id, req.Title, req.Content, req.Category, tags)
	if err != nil {
		writeNotFound(w, errPostNotFound)
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}

// DeletePost handles DELETE /posts/{id}.
func (h *PostsHandler) DeletePost(w http.ResponseWriter, r *http.Request, id generated.PostId) {
	err := h.Store.DeletePost(id)
	if err != nil {
		writeNotFound(w, errPostNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
