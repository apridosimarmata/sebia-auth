package auth

import (
	"context"
	"errors"
	"mini-wallet/domain"
	"mini-wallet/domain/user"
	"mini-wallet/utils"
	"net/http"
	"time"

	_auth "mini-wallet/domain/auth"
)

type authMiddleware struct {
	userRepository user.UserRepository
	config         *utils.AppConfig
}

func NewAuthMiddleware(repositories domain.Repositories, config *utils.AppConfig) _auth.AuthMiddleware {
	return &authMiddleware{
		userRepository: repositories.UserRepository,
		config:         config,
	}
}

func (middleware *authMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// processing access token
		tokenStatus, userId := middleware.processAccessToken(&w, r)
		if tokenStatus == _auth.ERROR_INVALID_TOKEN {
			// http.Error(w, "invalid access", http.StatusUnauthorized)
			// return
		}

		// refresh access
		if tokenStatus == _auth.ERROR_EXPIRED_TOKEN {
			err := middleware.refreshAccess(r.Context(), *userId, &w)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				// return
			}
			tokenStatus = 0
		}

		ctx := context.WithValue(r.Context(), _auth.UserIDContext{}, userId)
		ctx = context.WithValue(ctx, _auth.TokenStatus{}, tokenStatus)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *authMiddleware) processAccessToken(w *http.ResponseWriter, r *http.Request) (int, *string) {
	accessToken, err := r.Cookie(middleware.config.AccessTokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(*w, "No accessToken found", http.StatusUnauthorized)
			return _auth.ERROR_INVALID_TOKEN, nil
		}
		http.Error(*w, "Error retrieving accessToken", http.StatusInternalServerError)
		return _auth.ERROR_INVALID_TOKEN, nil
	}

	token := accessToken.Value
	userId, status := _auth.ExtractUserIDFromToken(token)
	if status != 0 {
		print(status)
		if status == _auth.ERROR_EXPIRED_TOKEN {
			return _auth.ERROR_EXPIRED_TOKEN, &userId
		}
		return _auth.ERROR_INVALID_TOKEN, nil
	}

	return 0, &userId
}

// writing  a new cookie of access token & refresh token to the response
func (middleware *authMiddleware) refreshAccess(ctx context.Context, userId string, w *http.ResponseWriter) error {
	user, err := middleware.userRepository.GetUserByUserID(ctx, userId)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	accessToken, _ := _auth.GenerateJWT(*user, "ACCESS")
	refreshToken, _ := _auth.GenerateJWT(*user, "REFRESH")

	now, _ := utils.GetJktTime()
	cookies := []*http.Cookie{
		{
			Name:     middleware.config.AccessTokenKey,
			Value:    accessToken,
			Domain:   ".sebia.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Hour * 24 * 31),
		},
		{
			Name:     middleware.config.RefreshTokenKey,
			Value:    refreshToken,
			Domain:   ".sebia.id",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  now.Add(time.Hour * 24 * 31),
		},
	}

	for _, cookie := range cookies {
		http.SetCookie(*w, cookie)
	}

	return nil
}
