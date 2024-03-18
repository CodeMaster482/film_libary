package middlware

import (
	"net/http"
	"runtime/debug"

	"films_library/pkg/logger"
)

type RecoveryMiddleware struct {
	log logger.Interface
}

func NewRecoveryMiddleware(log logger.Interface) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		log: log,
	}
}

func (m *RecoveryMiddleware) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				m.log.Fatal(rvr, debug.Stack())

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
