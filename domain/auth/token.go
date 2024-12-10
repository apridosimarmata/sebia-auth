package auth

import (
	"fmt"
	"mini-wallet/domain/user"
	"mini-wallet/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TOKEN_TYPE_REFRESH  = "REFRESH"
	ERROR_EXPIRED_TOKEN = 1
	ERROR_INVALID_TOKEN = 2
)

var secretKey = []byte("your-256-bit-secret")

func GenerateJWT(user user.UserEntity, tokenType string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	if tokenType == TOKEN_TYPE_REFRESH {
		expirationTime = time.Now().Add(30 * 24 * time.Hour)
	}

	claims := AcessTokenClaims{
		Name:   user.Name,
		UserID: user.UID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://dev.sebia.id",
			Subject:   user.UID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractUserIDFromToken(tokenString string) (string, int) {
	// Parse the token
	token, _ := jwt.ParseWithClaims(tokenString, &AcessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Validate the token and check the claims
	if token == nil || token.Claims.(*AcessTokenClaims) == nil {
		return "", ERROR_INVALID_TOKEN
	}

	if claims, ok := token.Claims.(*AcessTokenClaims); ok {
		now, _ := utils.GetJktTime()
		if claims.ExpiresAt.Time.Before(*now) {
			return claims.Subject, ERROR_EXPIRED_TOKEN
		}
		return claims.Subject, 0
	} else {
		return "", ERROR_INVALID_TOKEN
	}
}

func ValidateToken(tokenString string) (*AcessTokenClaims, int) {
	token, _ := jwt.ParseWithClaims(tokenString, &AcessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if token == nil || token.Claims.(*AcessTokenClaims) == nil {
		return nil, ERROR_INVALID_TOKEN
	}

	if claims, ok := token.Claims.(*AcessTokenClaims); ok {
		now, _ := utils.GetJktTime()
		if claims.ExpiresAt.Time.Before(*now) {
			return claims, ERROR_EXPIRED_TOKEN
		}
		return claims, 0
	}

	return nil, ERROR_INVALID_TOKEN
}
