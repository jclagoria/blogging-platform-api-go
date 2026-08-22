package handlers

import "github.com/juanka/blogging-platform-api/internal/generated"

// PostStore defines the data access contract for blog posts.
type PostStore interface {
	CreatePost(title, content, category string, tags []string) (*generated.Post, error)
	GetPost(id int) (*generated.Post, error)
	ListPosts(term string) ([]generated.Post, error)
	UpdatePost(id int, title, content, category string, tags []string) (*generated.Post, error)
	DeletePost(id int) error
}
