package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/juanka/blogging-platform-api/internal/generated"
)

// InMemoryPostStore is a PostStore backed by an in-memory map.
// ponytail: in-memory store for dev/test, swap for MongoDB in production
type InMemoryPostStore struct {
	mu      sync.RWMutex
	posts   map[int]*generated.Post
	nextID  int
}

func NewInMemoryPostStore() *InMemoryPostStore {
	return &InMemoryPostStore{
		posts:  make(map[int]*generated.Post),
		nextID: 1,
	}
}

func (s *InMemoryPostStore) CreatePost(title, content, category string, tags []string) (*generated.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	post := &generated.Post{
		Id:        s.nextID,
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      &tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.posts[s.nextID] = post
	s.nextID++

	return post, nil
}

func (s *InMemoryPostStore) GetPost(id int) (*generated.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post %d not found", id)
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

func (s *InMemoryPostStore) UpdatePost(id int, title, content, category string, tags []string) (*generated.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post %d not found", id)
	}

	post.Title = title
	post.Content = content
	post.Category = category
	post.Tags = &tags
	post.UpdatedAt = time.Now().UTC()

	return post, nil
}

func (s *InMemoryPostStore) DeletePost(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return fmt.Errorf("post %d not found", id)
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
