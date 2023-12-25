package middlewares

import (
	"accommodations-service/security"
	"accommodations-service/utils"
	"log"
	"net/http"
)

func RoleValidator(ac *security.AccessControl, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role := ctx.Value("role")
		log.Println("ROLA JE ", role.(string))
		if !ac.IsAccessAllowed(role.(string), r.URL.Path) {
			log.Println("Path je,", r.URL.Path)
			utils.WriteErrorResp("Unathorized4", 401, r.URL.Path, w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
