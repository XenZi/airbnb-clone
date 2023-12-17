package middlewares

import (
	"auth-service/security"
	"auth-service/utils"
	"log"
	"net/http"
)


func RoleValidator(ac *security.AccessControl, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role := ctx.Value("role")
		log.Println(role.(string))
		if !ac.IsAccessAllowed(role.(string), r.URL.Path) {
			utils.WriteErrorResp("Unathorized", 401, r.URL.Path, w)
			return
		}
		next.ServeHTTP(w, r)	
	})
}
