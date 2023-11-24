package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r.Header.Get("Authorization"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized - Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			http.Error(w, "Unauthorized - User ID not found in token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)	
	})
}

func extractToken(authorizationHeader string) string {
	bearerToken := strings.Split(authorizationHeader, " ")

	if len(bearerToken) == 2 {
		return bearerToken[1]
	}

	return ""
}
