package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{}
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			response.Message = "missing authorization header"
			json.NewEncoder(w).Encode(response)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			w.WriteHeader(http.StatusUnauthorized)
			response.Message = "invalid token format"
			json.NewEncoder(w).Encode(response)
			return
		}

		claims, err := validateJWT(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response.Message = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			response.Message = "invalid token payload"
			json.NewEncoder(w).Encode(response)
			return
		}
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", int(userID)))

		next.ServeHTTP(w, r)
	})
}

func validateJWT(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte("myjwtsecret")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}
