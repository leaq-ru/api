package middleware

import (
	"github.com/rs/cors"
	"net/http"
)

var c = cors.New(cors.Options{
	AllowedMethods: []string{http.MethodDelete, http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodHead},
	AllowedHeaders: []string{"Authorization", "Content-Type"},
})

func CORS(next http.Handler) http.Handler {
	return c.Handler(next)
}
