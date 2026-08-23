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

const (
	mongoRegex   = "$regex"
	mongoOptions = "$options"
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

// postDoc is the BSON representation stored in MongoDB.
type postDoc struct {
	ID        bson.ObjectID `bson:"_id"`
	Title     string        `bson:"title"`
	Content   string        `bson:"content"`
	Category  string        `bson:"category"`
	Tags      []string      `bson:"tags,omitempty"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}

func (d *postDoc) toPost() *generated.Post {
	var tags *[]string
	if len(d.Tags) > 0 {
		tags = &d.Tags
	}
	return &generated.Post{
		Id:        d.ID.Hex(),
		Title:     d.Title,
		Content:   d.Content,
		Category:  d.Category,
		Tags:      tags,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (s *MongoPostStore) CreatePost(title, content, category string, tags []string) (*generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	doc := postDoc{
		ID:        bson.NewObjectID(),
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.col().InsertOne(ctx, doc); err != nil {
		return nil, fmt.Errorf("insert post: %w", err)
	}

	return doc.toPost(), nil
}

func (s *MongoPostStore) GetPost(id string) (*generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post id: %s", id)
	}

	var doc postDoc
	if err := s.col().FindOne(ctx, bson.M{"_id": oid}).Decode(&doc); err != nil {
		return nil, fmt.Errorf("post %s not found", id)
	}
	return doc.toPost(), nil
}

func (s *MongoPostStore) ListPosts(term string) ([]generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if term != "" {
		t := strings.ToLower(term)
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{mongoRegex: t, mongoOptions: "i"}},
				{"content": bson.M{mongoRegex: t, mongoOptions: "i"}},
				{"category": bson.M{mongoRegex: t, mongoOptions: "i"}},
			},
		}
	}

	cursor, err := s.col().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer func() { _ = cursor.Close(ctx) }()

	var docs []postDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode posts: %w", err)
	}

	posts := make([]generated.Post, len(docs))
	for i, doc := range docs {
		posts[i] = *doc.toPost()
	}
	return posts, nil
}

func (s *MongoPostStore) UpdatePost(id string, title, content, category string, tags []string) (*generated.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post id: %s", id)
	}

	update := bson.M{
		"$set": bson.M{
			"title":     title,
			"content":   content,
			"category":  category,
			"tags":      tags,
			"updatedAt": time.Now().UTC(),
		},
	}

	res, err := s.col().UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("post %s not found", id)
	}

	return s.GetPost(id)
}

func (s *MongoPostStore) DeletePost(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid post id: %s", id)
	}

	res, err := s.col().DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("post %s not found", id)
	}
	return nil
}
