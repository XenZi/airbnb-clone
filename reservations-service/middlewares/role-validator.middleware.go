package middlewares

import (
	"log"
	"net/http"
	"reservation-service/utils"
)

func RoleValidator(allowedRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role := ctx.Value("role")
		log.Println("Role is:", role.(string))

		// Check if the user has the required role

		if role.(string) != allowedRole {

			log.Println("Unauthorized. Required role:", allowedRole)
			utils.WriteErrorResp("Unauthorized", 401, r.URL.Path, w)
			return
		}
		log.Println("USLO JE U ROLE VALIDATOR")

		next.ServeHTTP(w, r)
	})
}
