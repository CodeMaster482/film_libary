package middlware

import (
	"net/http"
	"strings"

	"films_library/pkg/response"
)

func Authentication(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session_id")
		if err != nil || cookie == nil {
			response.ErrorResponse(w, http.StatusUnauthorized, "missing token unauthorized", nil)
			return
		}

		if cookie.Value != "token_admin" || cookie.Name != "token_user" {
			response.ErrorResponse(w, http.StatusUnauthorized, "invalid token", nil)
			return
		}

		if cookie.Value == "token_user" {
			if r.Method != http.MethodGet {
				response.ErrorResponse(w, http.StatusForbidden, "forbidden", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})

	return http.HandlerFunc(fn)
}
