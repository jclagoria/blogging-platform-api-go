package handlers

import (
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/juanka/blogging-platform-api/internal/generated"
)

// InMemoryPostStore is a PostStore backed by an in-memory map.
// ponytail: in-memory store for dev/test, swap for MongoDB in production
type InMemoryPostStore struct {
	mu    sync.RWMutex
	posts map[string]*generated.Post
}

func NewInMemoryPostStore() *InMemoryPostStore {
	return &InMemoryPostStore{
		posts: make(map[string]*generated.Post),
	}
}

func generateID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *InMemoryPostStore) CreatePost(title, content, category string, tags []string) (*generated.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	id := generateID()
	post := &generated.Post{
		Id:        id,
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      &tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.posts[id] = post

	return post, nil
}

func (s *InMemoryPostStore) GetPost(id string) (*generated.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post %s not found", id)
	}
	return post, nil
}

func (s *InMemoryPostStore) ListPosts(term string) ([]generated.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []generated.Post
	for _, p := range s.posts {
		if term == "" || matchesTerm(p, term) {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (s *InMemoryPostStore) UpdatePost(id string, title, content, category string, tags []string) (*generated.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post %s not found", id)
	}

	post.Title = title
	post.Content = content
	post.Category = category
	post.Tags = &tags
	post.UpdatedAt = time.Now().UTC()

	return post, nil
}

func (s *InMemoryPostStore) DeletePost(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return fmt.Errorf("post %s not found", id)
	}
	delete(s.posts, id)
	return nil
}

func matchesTerm(p *generated.Post, term string) bool {
	t := strings.ToLower(term)
	return strings.Contains(strings.ToLower(p.Title), t) ||
		strings.Contains(strings.ToLower(p.Content), t) ||
		strings.Contains(strings.ToLower(p.Category), t)
}
