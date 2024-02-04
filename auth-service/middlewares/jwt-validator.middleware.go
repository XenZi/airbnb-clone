package middlewares

import (
	"auth-service/utils"
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
			utils.WriteErrorResp("Unathorized", 401, r.URL.Path, w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteErrorResp("Unathorized", 401, r.URL.Path, w)
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			utils.WriteErrorResp("Unathorized", 401, r.URL.Path, w)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "role", claims["role"])
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
