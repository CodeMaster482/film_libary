package middlware

import (
	"films_library/pkg/response"
	"net/http"
)

func AllowedMethod(next http.Handler) http.Handler {
	allowedEndpoints := map[string]map[string]bool{
		"/actor": {
			"GET": true,
		},
		"/actor/add": {
			"POST": true,
		},
		"/actor/update": {
			"PUT": true,
		},
		"/actor/delete": {
			"DELETE": true,
		},
		"/film/add": {
			"POST": true,
		},
		"/film/update": {
			"PUT": true,
		},
		"/film": {
			"GET": true,
		},
		"/film/delete": {
			"DELETE": true,
		},
		"/film/search": {
			"GET": true,
		},
	}

	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedMethods, ok := allowedEndpoints[r.URL.Path]
		if !ok {
			response.ErrorResponse(w, http.StatusNotFound, "Not found", nil)
			return
		}

		if !allowedMethods[r.Method] {
			response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
			return
		}

		next.ServeHTTP(w, r)
	})

	return http.HandlerFunc(fn)
}
