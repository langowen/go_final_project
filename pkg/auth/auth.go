package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/langowen/go_final_project/pkg/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrNoPassword   = errors.New("password not set")
)

type Claims struct {
	jwt.RegisteredClaims
	PasswordHash string `json:"pwd_hash"`
}

func GenerateToken(cfg *config.Config) (string, error) {
	if cfg.Token == "" {
		return "", ErrNoPassword
	}

	passwordHash := hashPassword(cfg.Token)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
		PasswordHash: passwordHash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Token))
}

func ValidateToken(cfg *config.Config, tokenString string) (bool, error) {
	if cfg.Token == "" {
		return false, ErrNoPassword
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Token), nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		currentHash := hashPassword(cfg.Token)
		return claims.PasswordHash == currentHash, nil
	}

	return false, ErrInvalidToken
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func Middleware(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.Token == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		valid, err := ValidateToken(cfg, cookie.Value)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
