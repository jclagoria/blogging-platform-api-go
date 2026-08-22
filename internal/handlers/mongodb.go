package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juanka/blogging-platform-api/internal/generated"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoPostStore is a PostStore backed by MongoDB.
type MongoPostStore struct {
	client     *mongo.Client
	database   string
	collection string
}

// NewMongoPostStore connects to MongoDB and returns a PostStore.
func NewMongoPostStore(uri, dbName string) (*MongoPostStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	return &MongoPostStore{
		client:     client,
		database:   dbName,
		collection: "posts",
	}, nil
}

// Close disconnects from MongoDB.
func (s *MongoPostStore) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *MongoPostStore) col() *mongo.Collection {
	return s.client.Database(s.database).Collection(s.collection)
}

func (s *MongoPostStore) CreatePost(title, content, category string, tags []string) (*generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	doc := bson.M{
		"title":     title,
		"content":   content,
		"category":  category,
		"tags":      tags,
		"createdAt": now,
		"updatedAt": now,
	}

	res, err := s.col().InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("insert post: %w", err)
	}

	post := &generated.Post{
		Id:        int(res.InsertedID.(bson.ObjectID).Timestamp().Unix()),
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      &tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return post, nil
}

func (s *MongoPostStore) GetPost(id int) (*generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var post generated.Post
	if err := s.col().FindOne(ctx, bson.M{"_id": id}).Decode(&post); err != nil {
		return nil, fmt.Errorf("post %d not found", id)
	}
	return &post, nil
}

func (s *MongoPostStore) ListPosts(term string) ([]generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if term != "" {
		t := strings.ToLower(term)
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": t, "$options": "i"}},
				{"content": bson.M{"$regex": t, "$options": "i"}},
				{"category": bson.M{"$regex": t, "$options": "i"}},
			},
		}
	}

	cursor, err := s.col().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	var posts []generated.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("decode posts: %w", err)
	}
	return posts, nil
}

func (s *MongoPostStore) UpdatePost(id int, title, content, category string, tags []string) (*generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":     title,
			"content":   content,
			"category":  category,
			"tags":      tags,
			"updatedAt": time.Now().UTC(),
		},
	}

	res, err := s.col().UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("post %d not found", id)
	}

	return s.GetPost(id)
}

func (s *MongoPostStore) DeletePost(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.col().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("post %d not found", id)
	}
	return nil
}
