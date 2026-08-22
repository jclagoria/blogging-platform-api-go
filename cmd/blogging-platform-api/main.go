package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/juanka/blogging-platform-api/internal/docs"
	"github.com/juanka/blogging-platform-api/internal/generated"
	"github.com/juanka/blogging-platform-api/internal/handlers"
	"github.com/juanka/blogging-platform-api/internal/middleware"
)

func main() {
	_ = godotenv.Load() // load .env if present, ignore error if missing

	r := chi.NewRouter()

	r.Use(middleware.CORS())

	var store handlers.PostStore
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		dbName := os.Getenv("MONGODB_DATABASE")
		if dbName == "" {
			dbName = "blogging_platform"
		}
		mongoStore, err := handlers.NewMongoPostStore(uri, dbName)
		if err != nil {
			log.Fatalf("MongoDB connection failed: %v", err)
		}
		defer func() { _ = mongoStore.Close(context.Background()) }()
		store = mongoStore
		log.Println("Connected to MongoDB")
	} else {
		store = handlers.NewInMemoryPostStore()
		log.Println("Using in-memory store (set MONGODB_URI for MongoDB)")
	}

	handler := &handlers.PostsHandler{Store: store}
	generated.HandlerFromMux(handler, r)

	apiDocs, _ := fs.Sub(docs.StaticFiles, "api")
	r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.FS(apiDocs))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting server on :%s", port)
		log.Printf("Swagger UI: http://localhost:%s/docs/swagger-ui/", port)
		log.Printf("Redoc:      http://localhost:%s/docs/redoc/", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
