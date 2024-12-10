package auth

import (
	"context"
	_auth "mini-wallet/domain/auth"
	"net/http"
)

func (middleware *authMiddleware) PublicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, userId := middleware.getAccessTokenUserId(&w, r) // return empty string if no user id found

		ctx := context.WithValue(r.Context(), _auth.UserIDContext{}, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *authMiddleware) getAccessTokenUserId(w *http.ResponseWriter, r *http.Request) (int, *string) {
	accessToken, err := r.Cookie(middleware.config.AccessTokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			return _auth.ERROR_INVALID_TOKEN, nil
		}
		return _auth.ERROR_INVALID_TOKEN, nil
	}

	token := accessToken.Value
	userId, _ := _auth.ExtractUserIDFromToken(token) // ignore status
	return 0, &userId
}
