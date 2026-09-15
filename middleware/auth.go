package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"ticket-system/utils"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONError(w, http.StatusUnauthorized, "authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				utils.JSONError(w, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := parts[1]

			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				utils.JSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			objectID, err := primitive.ObjectIDFromHex(claims.UserID)
			if err != nil {
				utils.JSONError(w, http.StatusUnauthorized, "invalid user id in token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, objectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (primitive.ObjectID, bool) {
	val := ctx.Value(UserIDContextKey)
	if val == nil {
		return primitive.NilObjectID, false
	}
	id, ok := val.(primitive.ObjectID)
	return id, ok
}
