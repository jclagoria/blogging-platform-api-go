package handlers

import "github.com/juanka/blogging-platform-api/internal/generated"

// PostStore defines the data access contract for blog posts.
type PostStore interface {
	CreatePost(title, content, category string, tags []string) (*generated.Post, error)
	GetPost(id string) (*generated.Post, error)
	ListPosts(term string) ([]generated.Post, error)
	UpdatePost(id string, title, content, category string, tags []string) (*generated.Post, error)
	DeletePost(id string) error
}
