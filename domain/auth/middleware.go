package auth

import (
	"net/http"
)

type AuthMiddleware interface {
	AuthMiddleware(next http.Handler) http.Handler
	PublicMiddleware(next http.Handler) http.Handler
}
