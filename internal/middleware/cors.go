package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORS returns a middleware that handles CORS headers.
func CORS() func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	})
	return c.Handler
}
