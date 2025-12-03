package mdware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxkey string

const ContextUserId ctxkey = "userId"

var jwtSecret = []byte(os.Getenv("JWTKEY"))

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}
			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalido token", http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(jwt.MapClaims)
			sub := claims["sub"]
			var userID uint
			switch v := sub.(type) {
			case float64:
				userID = uint(v)
			case float32:
				userID = uint(v)
			case int:
				userID = uint(v)
			case int64:
				userID = uint(v)
			default:
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserId, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
